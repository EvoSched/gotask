package handler

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/EvoSched/gotask/internal/models"

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
			t, err := parseTask(args)
			if err != nil {
				log.Fatal(err)
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
			fmt.Println("modifies task")
		},
	}
	return editCmd
}

func (h *Handler) NoteCmd() *cobra.Command {
	comCmd := &cobra.Command{
		Use:   "note",
		Short: "Notes a task by ID",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatal(errors.New("argument mismatch for comment command"))
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
			task, err := h.service.GetTask(id)
			if err != nil {
				return
			}
			fmt.Println(task)
			fmt.Println("Note: ", args[1])
		},
	}
	return comCmd
}

func (h *Handler) ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := h.service.GetTasks()
			if err != nil {
				fmt.Println("Error fetching tasks")
				log.Fatal(err)
				return
			}

			fmt.Println("Tasks:")
			for i, task := range tasks {
				str := "N/A"
				if task.TS != nil {
					str = task.TS.String()
				}
				fmt.Printf("%d. %s %s %s %d\n", i+1, task.Desc, str, task.Tags[0], task.Priority)
			}
		},
	}
	return listCmd
}

func (h *Handler) ComCmd() *cobra.Command {
	comCmd := &cobra.Command{
		Use:   "com",
		Short: "Complete a task by ID",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Complete task by ID")
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
		},
	}
	return incomCmd
}

func helper(s string) (*time.Time, *models.TimeStamp, error) {
	t, errD := parseDate(s)
	if errD != nil {
		t1, t2, errT := parseTimeStamp(s)
		if errT != nil {
			return nil, nil, errors.New("attempts to use invalid time statement")
		}
		ts := &models.TimeStamp{Start: t1, End: t2}
		return nil, ts, nil
	} else {
		return t, nil, nil
	}
}
