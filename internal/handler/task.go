package handler

import (
	"errors"
	"fmt"
	"github.com/EvoSched/gotask/internal/models"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

const (
	DateFmtDMY = "02-01-2006"
	DateFmtYMD = "2006-01-02"
	// todo need to add more date formats
	//		need to support hour, min, sec if the user desires
)

var dateFormats = []string{DateFmtDMY, DateFmtYMD}

func (h *Handler) RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "Task manager",
	}

	return rootCmd
}

func (h *Handler) AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal(errors.New("task name is required"))
			}
			t, err := parseAdd(args)
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

func (h *Handler) ComCmd() *cobra.Command {
	comCmd := &cobra.Command{
		Use:   "com",
		Short: "Comment a task by ID",
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
			fmt.Println("Comment: ", args[1])
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
				str := "nil"
				if task.Date != nil {
					str = task.Date.Format(dateFormats[0])
				}
				fmt.Printf("%d. %s %s %s %s\n", i+1, task.Title, task.Description, str, task.Tags[0])
			}
		},
	}

	return listCmd
}

/*
grammar:

gt add <title> <description> <day> at <time-time> <tag>

gt add <title> <description> <date> <time-time> <tag>

vs

gt add <title> <description> <due> <tag> <-- this is what we just finished
*/
func parseAdd(args []string) (*models.Task, error) {
	var description string
	var date time.Time
	var tags []string
	// we could in theory briefly check os.Args to see whether the first two arguments contain strings (for now, assume it does)
	title := args[0]
	if len(args) > 1 {
		description = args[1]
	}
	flg := false
	for _, arg := range args[2:] {
		if arg[0] == '+' {
			tags = append(tags, arg[1:])
		} else {
			if flg {
				return nil, fmt.Errorf("invalid command: contains repeat of time argument")
			}
			flg = true
			d, err := parseDate(arg)
			if err != nil {
				return nil, err
			}
			date = d
		}
	}

	return models.NewTask(0, title, description, &date, tags), nil
}

func parseDate(arg string) (time.Time, error) {
	today := time.Now()
	switch arg {
	case "eod":
		return today, nil
	case "eow", "sat":
		return today.AddDate(0, 0, int(time.Saturday-today.Weekday()+7)%7), nil
	case "sun":
		return today.AddDate(0, 0, int(time.Sunday-today.Weekday()+7)%7), nil
	case "mon":
		return today.AddDate(0, 0, int(time.Monday-today.Weekday()+7)%7), nil
	case "tue":
		return today.AddDate(0, 0, int(time.Tuesday-today.Weekday()+7)%7), nil
	case "wed":
		return today.AddDate(0, 0, int(time.Wednesday-today.Weekday()+7)%7), nil
	case "thu":
		return today.AddDate(0, 0, int(time.Thursday-today.Weekday()+7)%7), nil
	case "fri":
		return today.AddDate(0, 0, int(time.Friday-today.Weekday()+7)%7), nil
	default:
		for _, v := range dateFormats {
			t, err := time.Parse(v, arg)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid date format: %s", arg)
	}
}
