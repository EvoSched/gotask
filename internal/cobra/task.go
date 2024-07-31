package cobra

import (
	"fmt"
	"github.com/EvoSched/gotask/internal/types"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	DateFmtDMY = "02-01-2006"
)

var dateFormats = []string{DateFmtDMY, time.DateOnly}

func (c *Cmd) RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gt",
		Short: "GoTask",
		Long: `GoTask is a comprehensive cli application for managing your tasks both intuitively and efficiently. 
It allows you to add, list, mod, get, complete, import, export, and prioritize your tasks with ease.`,
	}
	return rootCmd
}

func (c *Cmd) AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Long: `Adds a new task with the provided description. Additional options include specifying time expression, tags, and priority.

Required:
- description  Description of the task to be added.

Optional:
- time      '@' marks the beginning of the time expression (halts when encountering non-time token).
- tag       Tag for categorizing the task, prefixed with '+'.
- priority  Priority level for the task from 1 to 10 (min-max), prefixed with '%'.`,
		Example: `gt add 'Write up ReadMe'
gt add 'Finish documentation' +work %8 @ 11-01-2024 10am-4:15
gt add "Setup database" @ 11-3 +project`,
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
			t := types.NewTask(*ti.desc, *ti.priority, ti.addTags, nil, ti.startAt, ti.endAt)
			i, err := c.repo.AddTask(t)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Added task %d.\n", i)
		},
	}

	return addCmd
}

func (c *Cmd) GetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get tasks by ID",
		Long: `Retrieves a task/s from the provided ids. Supports multiple retrievals in a single command.

Required:
- id  Id referencing task.`,
		Example: `gt get 1
gt get 1 3`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseGet(args)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range ids {
				t, err := c.repo.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				displayTask(t)
				fmt.Println()
			}
		},
	}
	return getCmd
}

func (c *Cmd) ModCmd() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "mod",
		Short: "Modify task by ID",
		Long: `Modifies an existing task given the arguments provided. Mod allows reorganization of arguments, but duplicates are invalid.

Optional:
- description  Description of the task to be modified. Must be surrounded by ' or " if description spans more than 1 word.
- time         '@' marks the beginning of the time expression (halts when encountering non-time token).
- tag          Tag for categorizing the task, prefixed with '+'.
- priority     Priority level for the task from 1 to 10 (min-max), prefixed with '%'.`,
		Example: `gt mod 1 'Reorganize structure of ReadMe'
gt mod 2 'Finish documentation for cobra commands' @ 11-01-2024 10am-4:15 +work %8
gt mod 3 +project "Setup database" @ 11-3`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ti, err := parseTask(args, false)
			if err != nil {
				log.Fatal(err)
			}
			t, err := c.repo.GetTask(*ti.id)
			if err != nil {
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
				t.EndAt = nil
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
			curr := time.Now()
			t.UpdatedAt = &curr
			err = c.repo.UpdateTask(t)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Update complete. 1 task modified.")
		},
	}
	return editCmd
}

func (c *Cmd) NoteCmd() *cobra.Command {
	comCmd := &cobra.Command{
		Use:   "note",
		Short: "Note task by ID",
		Long: `Attaches a note to a given task provided the task id. Notes are immutable and cannot be edited once created.

Required:
- id    Id referencing task.
- note  Note providing additional clarification for given task`,
		Example: `gt note 1 "Provide short gif demonstrating GoTask CLI and TUI"
gt note 2 "Finish writing up man docs from cobra commands"`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id, n, err := parseNote(args)
			if err != nil {
				log.Fatal(err)
			}
			d, err := c.repo.GetDesc(id)
			if err != nil {
				log.Fatal(err)
			}
			err = c.repo.AddNote(id, n)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Task %d '%s' has been updated with a new note:\n", id, d)
			fmt.Printf("  - Note: \"%s\"\n", n)
			fmt.Println("1 task updated with a note.")
		},
	}
	return comCmd
}

func (c *Cmd) ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all tasks",
		Long:    "Displays a list of all tasks created both new, overdue, and archived.",
		Example: "gt list",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			t, err := c.repo.GetTasks()
			if err != nil {
				log.Fatal(err)
			}
			displayTasks(t)
		},
	}
	return listCmd
}

func (c *Cmd) DueCmd() *cobra.Command {
	dueCmd := &cobra.Command{
		Use:     "due",
		Short:   "List all due and overdue tasks",
		Long:    "Displays a list of all tasks due and overdue.",
		Example: "gt due",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			t, err := c.repo.GetTasksDue()
			if err != nil {
				log.Fatal(err)
			}
			displayDueTasks(t)
		},
	}
	return dueCmd
}

func (c *Cmd) ArchivedCmd() *cobra.Command {
	dueCmd := &cobra.Command{
		Use:     "archived",
		Short:   "List all archived tasks",
		Long:    "Displays a list of all tasks archived.",
		Example: "gt archived",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			t, err := c.repo.GetTasksArchived()
			if err != nil {
				log.Fatal(err)
			}
			displayArchivedTasks(t)
		},
	}
	return dueCmd
}

func (c *Cmd) DoneCmd() *cobra.Command {
	doneCmd := &cobra.Command{
		Use:     "done",
		Short:   "Mark task as complete by ID",
		Long:    "Marks all tasks provided by ID as complete. This updates the lists that the tasks will now appear in (e.g. due, archived)",
		Example: "gt done 2\ngt done 1 3",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseDone(args)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range ids {
				t, err := c.repo.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				if t.Finished {
					fmt.Printf("Task %d already finished.\n", i)
					return
				} else {
					err = c.repo.UpdateStatus(i, true)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Finished task %d '%s'.\n", i, t.Desc)
				}
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

func (c *Cmd) UndoCmd() *cobra.Command {
	undoCmd := &cobra.Command{
		Use:     "undo",
		Short:   "Mark task as incomplete by ID",
		Long:    "Marks all tasks provided by ID as incomplete. This updates the lists that the tasks will now appear in (e.g. due, archived)",
		Example: "gt undo 3\ngt undo 2 1",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseDone(args)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range ids {
				t, err := c.repo.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				if !t.Finished {
					fmt.Printf("Task %d already incomplete.\n", i)
					return
				} else {
					err = c.repo.UpdateStatus(i, false)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Reverted task %d '%s' to incomplete.\n", i, t.Desc)
				}
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

func (c *Cmd) DeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete tasks by ID",
		Long:    "Deletes all tasks provided by ID. This updates the lists that the tasks will no longer appear in (e.g. due, archived, list)",
		Example: "gt delete 1\ngt delete 2 3",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := parseGet(args)
			if err != nil {
				log.Fatal(err)
			}
			var tasks []*types.Task
			for _, i := range ids {
				t, err := c.repo.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				tasks = append(tasks, t)
			}
			fmt.Printf("Preparing to delete tasks with ")
			if len(ids) == 1 {
				fmt.Printf("ID: %d\n", ids[0])
			} else {
				fmt.Printf("IDs: %d", ids[0])
				for j := 1; j < len(ids); j++ {
					fmt.Printf(", %d", ids[j])
				}
			}
			fmt.Println()
			for _, t := range tasks {
				fmt.Printf("  - Task %d: '%s'\n", t.ID, t.Desc)
			}
			if len(ids) == 1 {
				fmt.Printf("\nAre you sure you want to delete this task? (y/n): ")
			} else {
				fmt.Printf("\nAre you sure you want to delete these tasks? (y/n): ")
			}
			var s string
			_, err = fmt.Scanln(&s)
			if err != nil {
				return
			}
			s = strings.TrimSpace(s)
			if len(s) > 0 && s[0] == 'y' || s[0] == 'Y' {
				for _, i := range ids {
					err = c.repo.DeleteTask(i)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	}
	return deleteCmd
}

//func (c *Cmd) ImportCmd() *cobra.Command {
//	importCmd := &cobra.Command{
//		Use:   "import",
//		Short: "Import tasks from a file",
//		Args:  cobra.ExactArgs(1),
//		Run: func(cmd *cobra.Command, args []string) {
//			fmt.Printf("Importing tasks from '%s'...\n", args[0])
//			fmt.Println("  - 10 tasks found in the file\n  - 9 tasks successfully imported\n  - 1 task skipped (duplicate ID: 5)\n\nImport complete. 9 tasks added.")
//		},
//	}
//	return importCmd
//}
//
//func (c *Cmd) ExportCmd() *cobra.Command {
//	exportCmd := &cobra.Command{
//		Use:   "export",
//		Short: "Export tasks to file",
//		Args:  cobra.ExactArgs(1),
//		Run: func(cmd *cobra.Command, args []string) {
//			fmt.Printf("Exporting tasks to '%s'...\n", args[0])
//			fmt.Printf("  - 15 tasks exported\n\n")
//			fmt.Printf("Export complete. All tasks saved to '%s'.\n", args[0])
//		},
//	}
//	return exportCmd
//}

func displayTask(task *types.Task) {
	// Print header
	fmt.Println("Task Details:")
	fmt.Println("--------------")
	fmt.Printf("ID             %d\n", task.ID)
	fmt.Printf("Description    %s\n", task.Desc)
	fmt.Printf("Priority       %d\n", task.Priority)
	if len(task.Tags) > 0 {
		t := strings.Join(task.Tags, ", ")
		fmt.Printf("Tags           %v\n", t)
	}

	// Display 'Due' with date and time
	if task.StartAt == nil && task.EndAt == nil {
		fmt.Println("Due            <not set>")
	} else if task.StartAt != nil && task.EndAt != nil {
		fmt.Printf("Due            %s %s - %s\n", task.StartAt.Format("Mon, 02 Jan 2006"), task.StartAt.Format(time.Kitchen), task.EndAt.Format(time.Kitchen))
	} else if task.StartAt != nil {
		fmt.Printf("Due            %s %s\n", task.StartAt.Format("Mon, 02 Jan 2006"), task.StartAt.Format(time.Kitchen))
	} else {
		fmt.Printf("Due            %s %s\n", task.EndAt.Format("Mon, 02 Jan 2006"), task.EndAt.Format(time.Kitchen))
	}

	// Display last modified time
	fmt.Printf("Last modified  %s\n", task.UpdatedAt.Format(time.RFC1123))

	fmt.Printf("\nNotes:\n")
	if task.Notes != nil {
		for _, n := range task.Notes {
			fmt.Printf("  - %s\n", n)
		}
	}
}

// formatTask prints a task in the desired format
func formatTask(task *types.Task) string {
	// Format the status
	status := "[ ]"
	if task.Finished {
		status = "[x]"
	}

	var d string
	if len(task.Desc) > 27 {
		d = task.Desc[:27]
		d += ".."
	} else {
		d = task.Desc
	}

	// Format the tags
	tags := strings.Join(task.Tags, ", ")
	if len(tags) > 10 {
		tags = tags[:10]
		tags += ".."
	}

	// Format the due date and time
	var due string
	if task.StartAt != nil && task.EndAt != nil {
		due = fmt.Sprintf("%s %s - %s",
			task.StartAt.Format("Mon, 02 Jan 2006"),
			task.StartAt.Format("03:04pm"),
			task.EndAt.Format("03:04pm"))
	} else if task.StartAt != nil {
		due = fmt.Sprintf("%s %s",
			task.StartAt.Format("Mon, 02 Jan 2006"),
			task.StartAt.Format("03:04pm"))
	} else {
		due = "-"
	}
	// Format the output string with additional spaces for the 'Due' column
	return fmt.Sprintf("%-6d %-7s %-30s %-10d %-13s %s   ", // Adjusted format string with extra spaces
		task.ID, status, d, task.Priority, tags, due)
}

// displayTasks prints a list of tasks in the desired format
func displayTasks(tasks []*types.Task) {
	// Print header
	fmt.Println("ID     Status  Desc                           Priority   Tags          Due   ")
	fmt.Println("---------------------------------------------------------------------------------------------------------")

	// Print each task
	for _, task := range tasks {
		fmt.Println(formatTask(task))
	}
}

func formatTaskArchived(task *types.Task, archived bool) string {
	var d string
	if len(task.Desc) > 27 {
		d = task.Desc[:27]
		d += ".."
	} else {
		d = task.Desc
	}

	// Format the tags
	tags := strings.Join(task.Tags, ", ")
	if len(tags) > 10 {
		tags = tags[:10]
		tags += ".."
	}

	var due string
	if !archived {
		if task.StartAt != nil && task.EndAt != nil {
			due = fmt.Sprintf("%s %s - %s",
				task.StartAt.Format("Mon, 02 Jan 2006"),
				task.StartAt.Format("03:04pm"),
				task.EndAt.Format("03:04pm"))
		} else if task.StartAt != nil {
			due = fmt.Sprintf("%s %s",
				task.StartAt.Format("Mon, 02 Jan 2006"),
				task.StartAt.Format("03:04pm"))
		} else {
			due = "-"
		}
	} else {
		due = task.CompletedAt.Format(time.DateTime)
	}

	// Format the output string with additional spaces for the 'Due' column
	return fmt.Sprintf("%-6d %-30s %-10d %-13s %s   ", // Adjusted format string with extra spaces
		task.ID, d, task.Priority, tags, due)
}

func displayDueTasks(tasks []*types.Task) {
	fmt.Println("ID     Desc                           Priority   Tags          Due   ")
	fmt.Println("-------------------------------------------------------------------------------------------------")

	for _, t := range tasks {
		fmt.Println(formatTaskArchived(t, false))
	}
}

func displayArchivedTasks(tasks []*types.Task) {
	fmt.Println("ID     Desc                           Priority   Tags          Completed   ")
	fmt.Println("-------------------------------------------------------------------------------------------------")

	for _, t := range tasks {
		fmt.Println(formatTaskArchived(t, true))
	}
}
