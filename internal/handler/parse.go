package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/EvoSched/gotask/internal/models"
)

type taskInfo struct {
	id       int
	desc     *string
	StartAt  *time.Time
	EndAt    *time.Time
	addTags  []string
	remTags  []string
	priority *int
}

type timeStamp struct {
	start *time.Time
	end   *time.Time
}

func parseTask(args []string, isAdd bool) (*taskInfo, error) {
	task := new(taskInfo)
	startIdx := 0
	if isAdd {
		task.desc = &args[0]
		startIdx++
	}

	var date *time.Time
	var tStmp *timeStamp
	timeFlg := false

	for i := startIdx; i < len(args); i++ {
		if args[i][0] == '+' {
			task.addTags = append(task.addTags, args[i][1:])
		} else if args[i][0] == '-' && !isAdd {
			task.remTags = append(task.remTags, args[i][1:])
		} else if !timeFlg && args[i][0] == '@' { // this requires that the time expression be separated from '@' ex. gt add "work" @ 12-3 +MA
			timeFlg = true
			c := i + 3
			j := i + 1
			for ; j < len(args) && j < c; j++ {
				t, ts, err := parseTime(args[j])
				if err != nil && date == nil && tStmp == nil {
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
					if tStmp != nil {
						return nil, errors.New("task timestamp already set")
					}
					tStmp = ts
				}
			}
			i = j - 1
		} else if task.priority == nil && args[i][0] == '%' {
			if len(args[i]) > 1 {
				p, err := strconv.Atoi(args[i][1:])
				if err != nil {
					return nil, err
				}
				task.priority = &p
			}
		} else if !isAdd && task.desc == nil {
			task.desc = &args[i]
		} else if isAdd {
			return nil, errors.New("attempts to use invalid prefix outside of valid set for add {+, @, %}")
		} else {
			return nil, errors.New("attempts to use invalid prefix outside of valid set for mod {+, -, @, %}")
		}
	}
	if date != nil && tStmp != nil {
		s := time.Date(date.Year(), date.Month(), date.Day(), tStmp.start.Hour(), tStmp.start.Minute(), 0, 0, time.UTC)
		if tStmp.end != nil {
			e := time.Date(date.Year(), date.Month(), date.Day(), tStmp.end.Hour(), tStmp.end.Minute(), 0, 0, time.UTC)
			task.StartAt = &s
			task.EndAt = &e
		} else {
			task.StartAt = &s
		}
	} else if date != nil {
		task.StartAt = date
	} else if tStmp != nil {
		task.StartAt = tStmp.start
		task.EndAt = tStmp.end
	}

	// default priority value
	if task.priority == nil {
		p := 5
		task.priority = &p
	}
	return task, nil
}

func parseList(args []string) ([]string, *timeStamp, error) {
	return nil, nil, nil
}

func parseGet(args []string) (*models.Task, error) {
	return nil, nil
}

func parseCom(args []string) ([]int, error) {
	var ids []int
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseNote(args []string) (int, string, error) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, "", errors.New("invalid number type entered for 'note' command")
	}
	return id, args[1], nil
}

func parseTime(s string) (*time.Time, *timeStamp, error) {
	t, errD := parseDate(s)
	if errD != nil {
		t1, t2, errT := parseTimeStamp(s)
		if errT != nil {
			return nil, nil, errors.New("attempts to use invalid time statement")
		}
		ts := &timeStamp{start: t1, end: t2}
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
		case 'a', 'A', 'p', 'P':
			if !hour || am {
				return nil, nil, fmt.Errorf("'am'/'pm' cannot occur without an hour nor can there be duplicates: %s", arg)
			}
			if i+1 < len(arg) && (arg[i+1] == 'm' || arg[i+1] == 'M') {
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
	if (startFormat == "pm" || startFormat == "PM") && startHour != 12 {
		startHour += 12
	} else if (startFormat == "am" || startFormat == "AM") && startHour == 12 {
		startHour = 0
	}

	if startMinute == -1 {
		startMinute = 0
	}

	// if we only have the first part of the timestamp
	if endHour == -1 {
		fmt.Printf("Start Time: %02d:%02d\n", startHour, startMinute)
		t := time.Now()
		st := time.Date(t.Year(), t.Month(), t.Day(), startHour, startMinute, 0, 0, t.Location())
		return &st, nil, nil
	}

	// Convert end time to 24-hour format
	if (endFormat == "pm" || endFormat == "PM") && endHour != 12 {
		endHour += 12
	} else if (endFormat == "am" || endFormat == "AM") && endHour == 12 {
		endHour = 0
	} else if endFormat == "" && endHour < startHour {
		// Handle cases where endFormat is not provided but endHour is in 12-hour format
		endHour += 12
	}

	if endMinute == -1 {
		endMinute = 0
	}

	if startHour > endHour || (startHour == endHour && startMinute >= endMinute) {
		return nil, nil, fmt.Errorf("starting time must be earlier than ending time: %s", arg)
	}

	fmt.Printf("Start Time: %02d:%02d\n", startHour, startMinute)
	fmt.Printf("End Time: %02d:%02d\n", endHour, endMinute)
	t := time.Now()
	st := time.Date(t.Year(), t.Month(), t.Day(), startHour, startMinute, 0, 0, t.Location())
	et := time.Date(t.Year(), t.Month(), t.Day(), endHour, endMinute, 0, 0, t.Location())
	return &st, &et, nil
}
