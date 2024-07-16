package handler

import (
	//TODO: sort imports
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
			//TODO: move validation or args to parseAdd function
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
				//TODO: i don't think we should display "nil" here. We should display an empty string, or "N/A" or something
				str := "nil"
				if task.Date != nil {
					str = task.Date.Format(dateFormats[0])
				}
				fmt.Printf("%d. %s %s %s\n", i+1, task.Description, str, task.Tags[0])
			}
		},
	}

	return listCmd
}

/*
grammar:

todo gt add <description> <day> at <time-time> <tag>

todo gt add <description> <date> <time-time> <tag>

gt add <description> <due> <tag> <-- this is what we just finished
*/
func parseAdd(args []string) (*models.Task, error) {
	//TODO: we should validate args here and return an error if they are invalid

	var description string
	var date *time.Time
	var tags []string
	// we could in theory briefly check os.Args to see whether the first two arguments contain strings (for now, assume it does)
	if len(args) > 0 {
		description = args[0]
	} else {
		return nil, errors.New("task name is required")
	}

	if len(args) > 1 {
		flg := false
		for _, arg := range args[1:] {
			if arg[0] == '+' {
				tags = append(tags, arg[1:])
			} else {
				if flg {
					return nil, errors.New("invalid command: contains repeat of time argument")
				}
				flg = true
				d, err := parseDate(arg)
				if err != nil {
					return nil, err
				}
				date = d
			}
		}
	}

	return models.NewTask(0, description, date, tags), nil
}

func parseDate(arg string) (*time.Time, error) {
	today := time.Now()
	switch arg {
	case "eod":
		return &today, nil
	case "eow", "sat":
		d := today.AddDate(0, 0, int(time.Saturday-today.Weekday()+7)%7)
		return &d, nil
	case "sun":
		d := today.AddDate(0, 0, int(time.Sunday-today.Weekday()+7)%7)
		return &d, nil
	case "mon":
		d := today.AddDate(0, 0, int(time.Monday-today.Weekday()+7)%7)
		return &d, nil
	case "tue":
		d := today.AddDate(0, 0, int(time.Tuesday-today.Weekday()+7)%7)
		return &d, nil
	case "wed":
		d := today.AddDate(0, 0, int(time.Wednesday-today.Weekday()+7)%7)
		return &d, nil
	case "thu":
		d := today.AddDate(0, 0, int(time.Thursday-today.Weekday()+7)%7)
		return &d, nil
	case "fri":
		d := today.AddDate(0, 0, int(time.Friday-today.Weekday()+7)%7)
		return &d, nil
	default:
		for _, v := range dateFormats {
			t, err := time.Parse(v, arg)
			if err == nil {
				return &t, nil
			}
		}
		return nil, fmt.Errorf("invalid date format: %s", arg)
	}
}

func parseTimeStamp(arg string) (*time.Time, *time.Time, error) {
	hour, colon, minute, am, dash := false, false, false, false, false
	startHour, endHour, startMinute, endMinute := 0, 0, 0, 0
	startFormat, endFormat := "", ""

	for index := 0; index < len(arg); {
		currentChar := string(arg[index])

		switch currentChar {
		case ":":
			if !hour || minute {
				return &time.Time{}, &time.Time{}, errors.New("colon cannot occur before hour or after minute")
			}

			if colon {
				return &time.Time{}, &time.Time{}, errors.New("colon cannot occur more than once in a given time")
			}

			colon = true

			index++
		case "-":
			if !hour || dash {
				return &time.Time{}, &time.Time{}, errors.New("dash cannot occur before hour or more than once in a given timestamp")
			}

			hour = false
			colon = false
			minute = false
			am = false
			dash = true

			index++
		default:
			if val, err := strconv.Atoi(currentChar); err == nil {
				if hour && !colon {
					return &time.Time{}, &time.Time{}, errors.New("hour cannot be more than two digits")
				}

				if minute {
					return &time.Time{}, &time.Time{}, errors.New("minute cannot be more than two digits")
				}

				if am {
					return &time.Time{}, &time.Time{}, errors.New("digit cannot occur directly after time format")
				}

				if index+1 >= len(arg) {
					return &time.Time{}, &time.Time{}, errors.New("invalid digit placement")
				}

				// If colon is true, then we're currently parsing a minute
				// If not, we're currently parsing an hour
				if colon {
					val_, err := strconv.Atoi(string(arg[index+1]))

					// If dash is true, we're parsing the ending half of the timestamp
					// If not, we're parsing the starting half
					if dash {
						if err != nil {
							return &time.Time{}, &time.Time{}, errors.New("minute must be more than two digits")
						} else {
							endMinute = (val * 10) + val_
							minute = true
							index += 2
						}
					} else {
						if err != nil {
							return &time.Time{}, &time.Time{}, errors.New("minute must be more than two digits")
						} else {
							startMinute = (val * 10) + val_
							minute = true
							index += 2
						}
					}
				} else {
					if dash {
						endHour = (endHour * 10) + val
						hour = true
						index++
					} else {
						startHour = (startHour * 10) + val
						hour = true
						index++
					}
				}
			} else {
				if !hour {
					return &time.Time{}, &time.Time{}, errors.New("letter cannot occur before a digit in timestamp")
				}

				// "AM" must directly follow a minute
				if !minute {
					return &time.Time{}, &time.Time{}, errors.New("invalid letter placement in timestamp")
				}

				if index+1 >= len(arg) {
					return &time.Time{}, &time.Time{}, errors.New("invalid letter placement")
				}

				// If dash is true, we're parsing the ending half of the timestamp
				// If not, we're parsing the starting half
				if dash {
					endFormat = string(arg[index : index+2])
				} else {
					startFormat = string(arg[index : index+2])
				}

				am = true
				index += 2
			}
		}
	}

	// Convert from 12-hour to 24-hour clock format
	if startFormat == "pm" {
		startHour += 12
	}
	if endFormat == "pm" {
		endHour += 12
	}

	if endHour < startHour {
		return &time.Time{}, &time.Time{}, errors.New("ending hour must be greater than starting hour")
	}

	if startMinute < 0 || startMinute > 60 {
		return &time.Time{}, &time.Time{}, errors.New("minute must be a value between 0 and 60")
	}

	if endMinute < 0 || endMinute > 60 {
		return &time.Time{}, &time.Time{}, errors.New("minute must be a value between 0 and 60")
	}

	year, month, day := time.Now().Date()
	startTime := time.Date(year, month, day, startHour, startMinute, 0, 0, time.UTC)
	endTime := time.Date(year, month, day, endHour, endMinute, 0, 0, time.UTC)

	return &startTime, &endTime, nil
}
