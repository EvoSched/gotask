package handler

import (
	"errors"
	"fmt"
	"github.com/EvoSched/gotask/internal/models"
	"log"
	"strconv"
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
			t := new(models.Task)
			t.ID = 0
			t.Desc = *ti.desc
			if ti.priority != nil {
				t.Priority = *ti.priority
			}
			t.Tags = ti.addTags
			if ti.StartAt != nil {
				t.StartAt = ti.StartAt
			}
			if ti.EndAt != nil {
				t.EndAt = ti.EndAt
			}
			fmt.Println(t)
		},
	}

	return addCmd
}

func (h *Handler) GetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get tasks by ID",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal(errors.New("task id is required"))
			}
			var ids []int
			for _, i := range args {
				i, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				ids = append(ids, i)
			}

			for _, i := range ids {
				t, err := h.service.GetTask(i)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Task:", t)
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
			//t, err := h.service.GetTask(ti.id)
			//if err != nil {
			//	return
			//}
			fmt.Println("Task:", ti)
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
			task, err := h.service.GetTask(id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(task)
			fmt.Println("Note:", n)
		},
	}
	return comCmd
}

func (h *Handler) ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			//tasks, err := h.service.GetTasks()
			//if err != nil {
			//	fmt.Println("Error fetching tasks")
			//	log.Fatal(err)
			//	return
			//}
			//
			//fmt.Println("Tasks:")
			//for i, task := range tasks {
			//	str := "N/A"
			//	if task.TS != nil {
			//		str = task.TS.String()
			//	}
			//	fmt.Printf("%d. %s %s %s %d\n", i+1, task.Desc, str, task.Tags[0], task.Priority)
			//}
		},
	}
	return listCmd
}

func (h *Handler) ComCmd() *cobra.Command {
	comCmd := &cobra.Command{
		Use:   "com",
		Short: "Complete a task by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Complete task by ID")
			ids, err := parseCom(args)
			if err != nil {
				return
			}
			fmt.Println(ids, "complete...")
		},
	}
	return comCmd
}

func (h *Handler) IncomCmd() *cobra.Command {
	incomCmd := &cobra.Command{
		Use:   "incom",
		Short: "Incomplete a task by ID",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Incomplete task by ID")
			ids, err := parseCom(args)
			if err != nil {
				return
			}
			fmt.Println(ids, "incomplete...")
		},
	}
	return incomCmd
}
