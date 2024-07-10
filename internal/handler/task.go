package handler

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type tAddArgs struct {
	ttl []string
	hr  []float32
	due []time.Time
	tag []string
}

const (
	DateFmtDMY = "02-01-2006"
	DateFmtYMD = "2006-01-02"
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
	command := []string{`"task"`, "4.5", "2000-01-01", `"hw3"`, "--MA", "--GoTask", "eod"}
	ta, err := parseAddTask(command)
	fmt.Println(ta, err)
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				log.Fatal(errors.New("not enough arguments"))
			}
			//command := []string{"add", `"task"`, "4.5", "2000-01-01"}
			//ta, err := parseAddTask(command)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return addCmd
}

func (h *Handler) GetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a task by ID",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Please provide a task ID")
				return
			}

			//transform id to int
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Error occurred")
				log.Fatal(err)
				return
			}

			//call service to get task
			task, err := h.service.GetTask(id)
			if err != nil {
				fmt.Println("Error fetching task")
				log.Fatal(err)
				return
			}

			fmt.Println("Task: ", task)
		},
	}

	return getCmd
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
				fmt.Printf("%d. %s\n", i+1, task.Title)
			}
		},
	}

	return listCmd
}

func parseAddTask(args []string) (*tAddArgs, error) {
	ta := new(tAddArgs)
	for _, arg := range args {
		if len(arg) > 3 && arg[0] == '"' && arg[len(arg)-1] == '"' {
			ta.ttl = append(ta.ttl, arg)
		} else if len(arg) > 2 && reflect.DeepEqual("--", arg[0:2]) {
			ta.tag = append(ta.tag, arg[2:])
		} else {
			x, err := strconv.ParseFloat(arg, 32)
			if err != nil {
				flg := false
				for _, layout := range dateFormats {
					t, err := time.Parse(layout, arg)
					if err == nil {
						ta.due = append(ta.due, t)
						flg = true
						break
					}
				}
				if !flg {
					t, err := parseTime(strings.ToLower(arg))
					if err != nil {
						return nil, err
					}
					ta.due = append(ta.due, t)
				}
			} else {
				ta.hr = append(ta.hr, float32(x))
			}
		}
	}
	return ta, nil
}

func parseTime(arg string) (time.Time, error) {
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
		return time.Time{}, fmt.Errorf("invalid date format: %s", arg)
	}
}
