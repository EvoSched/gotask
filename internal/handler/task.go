package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

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
			prompt := promptui.Prompt{
				Label: "Enter Task",
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			//call service to create task

			fmt.Printf("Added task: %s\n", result)
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

				if index + 1 >= len(arg) {
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
						}
					} else {
						if err != nil {
							return &time.Time{}, &time.Time{}, errors.New("minute must be more than two digits")
						} else {
							startMinute = (val * 10) + val_
						}
					}

					index += 2

					minute = true
				} else {
					val_, err := strconv.Atoi(string(arg[index+1]))

					if dash {
						if err != nil {
							endHour = val
							index++
						} else {
							endHour = (val * 10) + val_
							index += 2
						}
					} else {
						if err != nil {
							startHour = val
							index++
						} else {
							startHour = (val * 10) + val_
							index += 2
						}
					}

					hour = true
				}
			} else if strings.ToLower(currentChar) == "p" || strings.ToLower(currentChar) == "a" {
				if !hour || am {
					return &time.Time{}, &time.Time{}, errors.New("time format cannot occur before hour")
				}

				if index + 1 >= len(arg) || strings.ToLower(string(arg[index+1])) != "m" {
					return &time.Time{}, &time.Time{}, errors.New("invalid time format")
				}

				if dash {
					endFormat = currentChar + string(arg[index+1])
				} else {
					startFormat = currentChar + string(arg[index+1])
				}

				index += 2

				am = true
			} else {
				return &time.Time{}, &time.Time{}, errors.New("invalid timestamp")
			}
		}
	}

	if startMinute < 0 || startMinute > 59 ||
		endMinute < 0 || endMinute > 59 {
		return &time.Time{}, &time.Time{}, errors.New("minute cannot be less than 0 or more than 59")
	}

	if startHour < 0 || endHour < 0 {
		return &time.Time{}, &time.Time{}, errors.New("hour cannot be less than 0")
	}

	if startFormat != "" {
		if startHour > 12 {
			return &time.Time{}, &time.Time{}, errors.New("hour in 12-hour format cannot be higher than 12")
		}
	} else {
		if startHour > 23 {
			return &time.Time{}, &time.Time{}, errors.New("hour in 24-hour format cannot be higher than 23")
		}
	}

	if endFormat != "" {
		if endHour > 12 {
			return &time.Time{}, &time.Time{}, errors.New("hour in 12-hour format cannot be higher than 12")
		}
	} else {
		if endHour > 23 {
			return &time.Time{}, &time.Time{}, errors.New("hour in 24-hour format cannot be higher than 23")
		}
	}

	if startFormat == "am" && endFormat == "am" {
		if endHour < startHour || (endHour == startHour && endMinute <= startMinute) {
			return &time.Time{}, &time.Time{}, errors.New("end time cannot occur before start time")
		}
	}

	if endFormat == "" && startFormat == "" {
		if endHour < startHour || (endHour == startHour && endMinute < startMinute) {
			if startHour <= 12 {
				startFormat = "am"
			} else {
				startHour -= 12
				startFormat = "pm"
			}
			endFormat = "pm"
		}
	}

	if startFormat == "pm" && startHour != 12 {
		startHour += 12
	} else if startFormat == "am" && startHour == 12 {
		startHour = 0
	}

	if endFormat == "pm" && endHour != 12 {
		endHour += 12
	} else if endFormat == "am" && endHour == 12 {
		endHour = 0
	}

	start := time.Date(0, 0, 0, startHour, startMinute, 0, 0, time.UTC)
	end := time.Date(0, 0, 0, endHour, endMinute, 0, 0, time.UTC)

	if end.Before(start) {
		return &time.Time{}, &time.Time{}, errors.New("end time cannot occur before start time")
	}

	return &start, &end, nil
}