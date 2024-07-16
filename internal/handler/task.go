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
				if task.TS != nil {
					str = task.TS.String()
				}
				fmt.Printf("%d. %s %s %s %d\n", i+1, task.Desc, str, task.Tags[0], task.Priority)
			}
		},
	}

	return listCmd
}

func parseAdd(args []string) (*models.Task, error) {
	var description string
	var tags []string

	if len(args) > 0 {
		description = args[0]
	} else {
		return nil, errors.New("task name is required")
	}

	var date *time.Time
	var timeStamp *models.TimeStamp
	var priority *int

	timeFlg := false
	priorityFlg := false
	for i := 1; i < len(args); i++ {
		if args[i][0] == '+' {
			tags = append(tags, args[i][1:])
		} else if !timeFlg && args[i][0] == '@' { // this requires that the time expression be separated from '@' ex. gt add "work" @ 12-3 +MA
			timeFlg = true
			c := i + 3
			j := i + 1
			for ; j < len(args) && j < c; j++ {
				t, ts, err := helper(args[j])
				if err != nil && date == nil && timeStamp == nil {
					return nil, err
				} else if err != nil {
					continue
				}
				if t != nil {
					if date != nil {
						return nil, errors.New("task date already set")
					}
					date = t
				}
				if ts != nil {
					if timeStamp != nil {
						return nil, errors.New("task timestamp already set")
					}
					timeStamp = ts
				}
			}
			i = j - 1
		} else if !priorityFlg && args[i][0] == '%' {
			priorityFlg = true
			if len(args[i]) > 1 {
				p, err := strconv.Atoi(args[i][1:])
				if err != nil {
					return nil, err
				}
				priority = &p
			}
		} else {
			return nil, errors.New("attempts to use invalid prefix outside of valid set {+, @, %, #}")
		}
	}
	if date != nil && timeStamp != nil {
		s := time.Date(date.Year(), date.Month(), date.Day(), timeStamp.Start.Hour(), timeStamp.Start.Minute(), 0, 0, time.UTC)
		e := time.Date(date.Year(), date.Month(), date.Day(), timeStamp.End.Hour(), timeStamp.End.Minute(), 0, 0, time.UTC)
		timeStamp = &models.TimeStamp{
			Start: &s,
			End:   &e,
		}
	} else if date != nil {
		timeStamp = &models.TimeStamp{
			Start: date,
		}
	}

	// default priority value
	p := 5
	if priority != nil {
		p = *priority
	}

	return models.NewTask(0, description, timeStamp, tags, p), nil
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
	startHour, endHour, startMinute, endMinute := -1, -1, -1, -1
	startFormat, endFormat := "", ""

	for i := 0; i < len(arg); i++ {
		switch arg[i] {
		case ':':
			if !hour || minute {
				return nil, nil, fmt.Errorf("colon is expected after hour, not minute: %s", arg)
			} else if am {
				return nil, nil, fmt.Errorf("colon can never occur after 'am' or 'pm': %s", arg)
			} else if colon {
				return nil, nil, fmt.Errorf("colons cannot be duplicated for same hour, minute combination: %s", arg)
			}
			colon = true
		case '-':
			if !hour || dash {
				return nil, nil, fmt.Errorf("dashes require an hour and cannot be duplicated: %s", arg)
			}
			dash = true
			hour = false
			colon = false
			minute = false
			am = false
		case 'a', 'p':
			if !hour || am {
				return nil, nil, fmt.Errorf("'am'/'pm' cannot occur without an hour nor can there be duplicates: %s", arg)
			}
			if i+1 < len(arg) && arg[i+1] == 'm' {
				s := arg[i : i+2]
				if !dash {
					startFormat = s
				} else {
					endFormat = s
				}
			} else {
				return nil, nil, fmt.Errorf("provided time tag besides valid 'am'/'pm': %s", arg)
			}
			am = true
			i++
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if hour && !colon {
				return nil, nil, fmt.Errorf("minutes must be separated by colon from hours: %s", arg)
			} else if minute {
				return nil, nil, fmt.Errorf("minutes cannot be duplicated: %s", arg)
			} else if am {
				return nil, nil, fmt.Errorf("hours and minutes cannot come after 'am'/'pm' signature: %s", arg)
			} else if hour {
				minute = true
				if i+1 < len(arg) && arg[i+1] >= '0' && arg[i+1] <= '9' {
					x, _ := strconv.Atoi(arg[i : i+2])
					if x > 59 {
						return nil, nil, fmt.Errorf("minutes cannot be greater than 59: %s", arg)
					}
					if !dash {
						startMinute = x
					} else {
						endMinute = x
					}
					i++
				} else {
					return nil, nil, fmt.Errorf("minutes require 2 digits: %s", arg)
				}
			} else {
				hour = true
				if i+1 >= len(arg) || arg[i+1] < '0' || arg[i+1] > '9' {
					x, _ := strconv.Atoi(arg[i : i+1])
					if x > 12 {
						return nil, nil, fmt.Errorf("hours cannot be greater than 12: %s", arg)
					}
					if !dash {
						startHour = x
					} else {
						endHour = x
					}
				} else if i+1 < len(arg) && arg[i+1] >= '0' && arg[i+1] <= '9' {
					x, _ := strconv.Atoi(arg[i : i+2])
					if x > 12 {
						return nil, nil, fmt.Errorf("hours cannot be greater than 12: %s", arg)
					}
					if !dash {
						startHour = x
					} else {
						endHour = x
					}
					i++
				} else {
					return nil, nil, fmt.Errorf("minutes require 2 digits: %s", arg)
				}
			}
		default:
			return nil, nil, fmt.Errorf("invalid time format: %s", arg)
		}
	}

	// Convert start time to 24-hour format
	if startFormat == "pm" && startHour != 12 {
		startHour += 12
	} else if startFormat == "am" && startHour == 12 {
		startHour = 0
	}

	// Convert end time to 24-hour format
	if endFormat == "pm" && endHour != 12 {
		endHour += 12
	} else if endFormat == "am" && endHour == 12 {
		endHour = 0
	}

	// Handle cases where endFormat is not provided but endHour is in 12-hour format
	if endFormat == "" && endHour < startHour {
		endHour += 12
	}

	if startHour > endHour || (startHour == endHour && startMinute >= endMinute) {
		return nil, nil, fmt.Errorf("starting time must be earlier than ending time: %s", arg)
	}

	if startMinute == -1 {
		startMinute = 0
	}
	if endMinute == -1 {
		endMinute = 0
	}

	fmt.Printf("Start Time: %02d:%02d\n", startHour, startMinute)
	fmt.Printf("End Time: %02d:%02d\n", endHour, endMinute)
	t := time.Now()
	st := time.Date(t.Year(), t.Month(), t.Day(), startHour, startMinute, 0, 0, t.Location())
	et := time.Date(t.Year(), t.Month(), t.Day(), endHour, endMinute, 0, 0, t.Location())
	return &st, &et, nil
}
