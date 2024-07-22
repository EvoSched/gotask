package handler

import (
	"fmt"
	"github.com/EvoSched/gotask/internal/models"
	"log"
	"time"

	"github.com/spf13/cobra"
)

const (
	DateFmtDMY = "02-01-2006"
)

var dateFormats = []string{DateFmtDMY, time.DateOnly}

func (h *Handler) RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "GoTask",
		Long: `GoTask is a cli application for managing tasks efficiently. 
It allows you to add, list, mod, get, complete, and prioritize your tasks with ease.`,
	}
	return rootCmd
}

func (h *Handler) AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Long: `Adds a new task with the provided description. Additional options include specifying time expression, tags, and priority.

Required:
- description: Description of the task to be added.

Optional:
- time: '@' marks the beginning of the time expression (halts when encountering non-time token).
- tag: Tag for categorizing the task, prefixed with '+'.
- priority: Priority level for the task from 1 to 10 (min-max), prefixed with '%'.

Example usages:
gt add 'Write up ReadMe'
gt add 'Finish documentation' +work %8 @ 11-01-2024 10am-4:15
gt add "Setup database" @ 11-3 +project
`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ti, err := parseTask(args, true)
			if err != nil {
				log.Fatal(err)
			}
			if ti.priority == nil {
				p := 5
				ti.priority = &p
			}
			t := models.NewTask(1, *ti.desc, *ti.priority, ti.addTags, nil, ti.startAt, ti.endAt)
			fmt.Printf("Added task %d.\n", t.ID)
		},
	}

	return addCmd
}

func (h *Handler) GetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get tasks by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseGet(args)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range ids {
				t, err := h.service.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				//check if task exists, because of nil pointer dereference
				if t == nil {
					fmt.Printf("Task %d not found.\n", i)
					continue
				}
				models.DisplayTask(t)
				fmt.Println()
			}
		},
	}
	return getCmd
}

func (h *Handler) ModCmd() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "mod",
		Short: "Modify a task by ID",
		Run: func(cmd *cobra.Command, args []string) {
			ti, err := parseTask(args, false)
			if err != nil {
				log.Fatal(err)
			}
			t, err := h.service.GetTask(*ti.id)
			if err != nil {
				return
			}
			//check if task exists, because of nil pointer dereference
			if t == nil {
				fmt.Printf("Task %d not found.\n", *ti.id)
				return
			}
			fmt.Printf("Task %d '%s' has been updated:\n", t.ID, t.Desc)
			if ti.desc != nil {
				fmt.Printf("  - Description updated to '%s'\n", *ti.desc)
				t.Desc = *ti.desc
			}
			if ti.priority != nil {
				fmt.Printf("  - Priority updated from %d to %d\n", t.Priority, *ti.priority)
				t.Priority = *ti.priority
			}
			if ti.addTags != nil {
				for _, tg := range ti.addTags {
					fmt.Printf("  - Tag added: %s\n", tg)
				}
				t.Tags = append(t.Tags, ti.addTags...)
			}
			if ti.startAt != nil {
				t.StartAt = ti.startAt
			}
			if ti.endAt != nil {
				t.EndAt = ti.endAt
			}
			if ti.startAt != nil {
				fmt.Printf("  - Time updated to ")
				if t.StartAt != nil && t.EndAt != nil {
					fmt.Printf("%s %s - %s\n", t.StartAt.Format("Mon, 02 Jan 2006"), t.StartAt.Format(time.Kitchen), t.EndAt.Format(time.Kitchen))
				} else {
					fmt.Printf("%s %s\n", t.StartAt.Format("Mon, 02 Jan 2006"), t.StartAt.Format(time.Kitchen))
				}
			}
			fmt.Println("Update complete. 1 task modified.")
		},
	}
	return editCmd
}

func (h *Handler) NoteCmd() *cobra.Command {
	comCmd := &cobra.Command{
		Use:   "note",
		Short: "Notes a task by ID",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id, n, err := parseNote(args)
			if err != nil {
				log.Fatal(err)
			}
			t, err := h.service.GetTask(id)
			if err != nil {
				log.Fatal(err)
			}
			//check if task exists, because of nil pointer dereference
			if t == nil {
				fmt.Printf("Task %d not found.\n", id)
				return
			}
			fmt.Printf("Task %d '%s' has been updated with a new note:\n", t.ID, t.Desc)
			fmt.Printf("  - Note: \"%s\"\n", n)
			fmt.Println("1 task updated with a note.")
		},
	}
	return comCmd
}

func (h *Handler) ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			t, err := h.service.GetTasks()
			if err != nil {
				log.Fatal(err)
			}
			models.DisplayTasks(t)
		},
	}
	return listCmd
}

func (h *Handler) DoneCmd() *cobra.Command {
	doneCmd := &cobra.Command{
		Use:   "done",
		Short: "Marks task as complete by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseDone(args)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range ids {
				t, err := h.service.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				//check if task exists, because of nil pointer dereference
				if t == nil {
					fmt.Printf("Task %d not found.\n", i)
					continue
				}
				fmt.Printf("Finished task %d '%s'.\n", i, t.Desc)
			}
			if len(ids) > 1 {
				fmt.Printf("Finished %d tasks.\n", len(ids))
			} else {
				fmt.Printf("Finished 1 task.\n")
			}
		},
	}
	return doneCmd
}

func (h *Handler) UndoCmd() *cobra.Command {
	undoCmd := &cobra.Command{
		Use:   "undo",
		Short: "Marks task as incomplete by ID",
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseDone(args)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range ids {
				t, err := h.service.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				//check if task exists, because of nil pointer dereference
				if t == nil {
					fmt.Printf("Task %d not found.\n", i)
					continue
				}
				fmt.Printf("Reverted task %d '%s' to incomplete.\n", i, t.Desc)
			}
			if len(ids) > 1 {
				fmt.Printf("Reverted %d tasks.\n", len(ids))
			} else {
				fmt.Printf("Reverted 1 task.\n")
			}
		},
	}
	return undoCmd
}

func (h *Handler) DeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes tasks by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			/*
				Preparing to delete tasks with IDs: 1, 3, 5
				  - Task 1: 'finish some work'
				  - Task 3: 'study BSTs'
				  - Task 5: 'clean the house'

				Are you sure you want to delete these tasks? (y/n):
			*/
		},
	}
	return deleteCmd
}

func (h *Handler) ImportCmd() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import tasks from a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Importing tasks from '%s'...\n", args[0])
			fmt.Println("  - 10 tasks found in the file\n  - 9 tasks successfully imported\n  - 1 task skipped (duplicate ID: 5)\n\nImport complete. 9 tasks added.")
		},
	}
	return importCmd
}

func (h *Handler) ExportCmd() *cobra.Command {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export tasks to file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Exporting tasks to '%s'...\n", args[0])
			fmt.Printf("  - 15 tasks exported\n\n")
			fmt.Printf("Export complete. All tasks saved to '%s'.\n", args[0])
		},
	}
	return exportCmd
}
