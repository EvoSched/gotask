// Package cobra implements command-line parsing and execution for the gotask application
package cobra

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// taskInfo represents the structure of a task with all its properties
// All fields are pointers to allow for optional values
type taskInfo struct {
	id       *int       // Unique identifier for the task
	desc     *string    // Task description
	startAt  *time.Time // Start time of the task
	endAt    *time.Time // End time of the task
	addTags  []string   // Tags to be added to the task
	remTags  []string   // Tags to be removed from the task
	priority *int       // Task priority (1-5, where 1 is highest)
}

// timeStamp represents a time range with optional start and end times
type timeStamp struct {
	start *time.Time // Start time
	end   *time.Time // End time (optional)
}

// parseTask processes command line arguments to create or modify a task
// isAdd determines whether this is a new task (true) or modifying an existing task (false)
func parseTask(args []string, isAdd bool) (*taskInfo, error) {
	// Initialize a new taskInfo object
	task := new(taskInfo)

	// Handle new task creation
	if isAdd {
		// Set the task description from the first argument
		task.desc = &args[0]
	} else {
		// Handle existing task modification
		// Parse the task ID from the first argument
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		task.id = &id
	}

	// Initialize variables to track date and time parsing
	var date *time.Time
	var tStmp *timeStamp
	timeFlg := false

	// This loop processes each argument in the command after the first one
	// Example command: gt add "Complete homework" +school @tomorrow %1
	// args would be: ["Complete homework", "+school", "@", "tomorrow", "%1"]
	for i := 1; i < len(args); i++ {
		// CASE 1: Adding Tags
		// If argument starts with '+', it's a tag to add
		// Example: +school, +urgent, +work
		if args[i][0] == '+' {
			// args[i][1:] removes the '+' and takes the rest of the string
			// Example: "+school" becomes "school"
			task.addTags = append(task.addTags, args[i][1:])

			// CASE 2: Removing Tags (only for existing tasks)
			// If argument starts with '!' and we're modifying an existing task
			// Example: !school (removes 'school' tag)
		} else if args[i][0] == '!' && !isAdd {
			task.remTags = append(task.remTags, args[i][1:])

			// CASE 3: Setting Time/Date
			// If argument starts with '@' and we haven't processed time yet
			// The '@' must be separated from the time expression
			// Example: gt add "work" @ 12-3pm +MA
			//         NOT: gt add "work" @12-3pm +MA
		} else if !timeFlg && args[i][0] == '@' {
			timeFlg = true // Mark that we're processing time
			c := i + 3     // Look ahead up to 3 arguments for time components
			j := i + 1     // Start from next argument after '@'
			curIdx := 0    // Track how many time components we've processed

			// Look at up to 3 arguments after '@' for time/date information
			// This allows formats like: @ tomorrow 2pm
			//                          @ 2pm
			//                          @ 2-4pm
			for ; j < len(args) && j < c; j++ {
				curIdx++
				// Try to parse the argument as either a date or time
				t, ts, err := parseTime(args[j])

				// If parsing failed and we haven't found any valid time yet, return error
				if err != nil && date == nil && tStmp == nil {
					return nil, err
				} else if err != nil {
					// If parsing failed but we already have some time info, just skip this part
					continue
				}

				// If we got a date (like "tomorrow", "next week")
				if t != nil {
					// Can't set date twice
					if date != nil {
						return nil, errors.New("task date already set")
					}
					date = t
				}

				// If we got a timestamp (like "2pm", "2-4pm")
				if ts != nil {
					// Can't set timestamp twice
					if tStmp != nil {
						return nil, errors.New("task timestamp already set")
					}
					tStmp = ts
				}
			}
			// Skip the arguments we just processed
			i = j - 1

			// CASE 4: Setting Priority
			// If argument starts with '%' and priority isn't set yet
			// Example: %1 (highest priority) to %5 (lowest priority)
		} else if task.priority == nil && args[i][0] == '%' {
			// Check if there's a number after '%'
			if len(args[i]) > 1 {
				// Convert the priority string to number
				// Example: "%1" becomes 1
				p, err := strconv.Atoi(args[i][1:])
				if err != nil {
					return nil, err
				}
				task.priority = &p
			}

			// CASE 5: Updating Description (only for existing tasks)
			// If we're modifying a task and haven't set description yet
		} else if !isAdd && task.desc == nil {
			task.desc = &args[i]

			// CASE 6: Error Cases
			// If none of the above cases match, it's an invalid argument
		} else if isAdd {
			// For new tasks, only +, @, and % prefixes are allowed
			return nil, errors.New("attempts to use invalid prefix outside of valid set for add {+, @, %}")
		} else {
			// For existing tasks, only +, -, @, and % prefixes are allowed
			return nil, errors.New("attempts to use invalid prefix outside of valid set for mod {+, -, @, %}")
		}
	}

	// FINAL STEP: Combining Date and Time Information
	// We might have both a date (like "tomorrow") and a time (like "2-4pm")
	// We need to combine them correctly

	// CASE 1: We have both date and time
	// Example: @ tomorrow 2-4pm
	if date != nil && tStmp != nil {
		// Create a new datetime by combining:
		// - Date parts (year, month, day) from the date argument
		// - Time parts (hour, minute) from the time argument
		// Example: If date is "tomorrow" (2024-01-20) and time is "2:30pm"
		//         Result will be "2024-01-20 14:30:00"
		s := time.Date(
			date.Year(),          // Year from date (e.g., 2024)
			date.Month(),         // Month from date (e.g., January)
			date.Day(),           // Day from date (e.g., 20)
			tStmp.start.Hour(),   // Hour from time (e.g., 14 for 2pm)
			tStmp.start.Minute(), // Minute from time (e.g., 30)
			0,                    // Seconds (always 0)
			0,                    // Nanoseconds (always 0)
			time.UTC,             // Timezone
		)

		// If we have an end time (e.g., "2-4pm"), set both start and end
		if tStmp.end != nil {
			// Create end time similar to start time
			e := time.Date(
				date.Year(), date.Month(), date.Day(),
				tStmp.end.Hour(), tStmp.end.Minute(),
				0, 0, time.UTC,
			)
			task.startAt = &s
			task.endAt = &e
		} else {
			// If no end time, just set the start time
			task.startAt = &s
		}

		// CASE 2: We only have a date
		// Example: @ tomorrow
	} else if date != nil {
		task.startAt = date // Just use the date as is

		// CASE 3: We only have a time
		// Example: @ 2-4pm
	} else if tStmp != nil {
		task.startAt = tStmp.start // Set the start time
		task.endAt = tStmp.end     // Set the end time (might be nil)
	}

	return task, nil
}

// parseGet processes arguments for the 'get' command
// This function converts a list of string IDs into actual numbers
//
// Example usage:
//
//	gt get 1 2 3
//	args would be: ["1", "2", "3"]
//	returns: [1, 2, 3]
func parseGet(args []string) ([]int, error) {
	// Create an empty slice to store the converted numbers
	var ids []int

	// Loop through each argument
	for _, arg := range args {
		// Try to convert the string to a number
		// Example: "1" becomes 1
		id, err := strconv.Atoi(arg)

		// If conversion fails (e.g., "abc" can't become a number)
		// return an error
		if err != nil {
			return nil, err
		}

		// Add the converted number to our list
		ids = append(ids, id)
	}

	// Return the list of task IDs
	return ids, nil
}

// parseDone processes arguments for the 'done' command
// This function is similar to parseGet, but specifically for marking tasks as completed
//
// Example usage:
//
//	gt done 1 2 3    (mark tasks 1, 2, and 3 as completed)
//	gt done 5        (mark task 5 as completed)
//
// Input:
//
//	args: ["1", "2", "3"] (array of task IDs as strings)
//
// Output:
//   - Success: returns [1, 2, 3] (array of integers)
//   - Error: returns nil and error if any ID is not a valid number
func parseDone(args []string) ([]int, error) {
	// Create an empty slice to store the task IDs
	var ids []int

	// Process each argument one by one
	// Example: for args ["1", "2", "3"]:
	//   First iteration:  arg = "1"
	//   Second iteration: arg = "2"
	//   Third iteration:  arg = "3"
	for _, arg := range args {
		// Convert the string ID to a number
		// Example: "1" becomes 1
		id, err := strconv.Atoi(arg)

		// If conversion fails (e.g., "abc" is not a valid number)
		// return immediately with an error
		if err != nil {
			return nil, err
		}

		// Add the converted number to our list of IDs
		ids = append(ids, id)
	}

	// Return the complete list of task IDs to mark as done
	return ids, nil
}

// parseNote processes arguments for the 'note' command
// This function handles adding notes to existing tasks
//
// Example usage:
//
//	gt note 1 "Remember to include tests"
//	args would be: ["1", "Remember to include tests"]
//
// Input:
//
//	args[0]: Task ID as string (e.g., "1")
//	args[1]: The note content (e.g., "Remember to include tests")
//
// Output:
//   - Success: returns (taskID, noteContent, nil)
//   - Error: returns (0, "", error) if task ID is not a valid number
//
// Note: This function expects exactly 2 arguments:
//  1. The task ID
//  2. The note content
func parseNote(args []string) (int, string, error) {
	// Try to convert the first argument to a number (task ID)
	// Example: "1" becomes 1
	id, err := strconv.Atoi(args[0])

	// If the conversion fails (e.g., "abc" is not a valid number)
	// return an error with a helpful message
	if err != nil {
		return 0, "", errors.New("invalid number type entered for 'note' command")
	}

	// Return three values:
	// 1. The task ID (as a number)
	// 2. The note content (everything after the ID)
	// 3. nil (no error)
	return id, args[1], nil
}

// parseTime is the main time parsing function that handles both dates and times
// It tries to parse the input first as a date, then as a time if that fails
//
// Example inputs:
//   - Dates: "tomorrow", "mon", "2024-01-20"
//   - Times: "2pm", "2:30pm", "2pm-4pm"
//
// Returns:
//   - For dates: returns (dateTime, nil, nil)
//   - For times: returns (nil, timeStamp, nil)
//   - For errors: returns (nil, nil, error)
func parseTime(s string) (*time.Time, *timeStamp, error) {
	// First, try to parse as a date (tomorrow, mon, etc.)
	t, errD := parseDate(s)
	if errD != nil {
		// If it's not a date, try to parse as a time
		t1, t2, errT := parseTimeStamp(s)
		if errT != nil {
			// If both date and time parsing fail, return error
			return nil, nil, errors.New("attempts to use invalid time statement")
		}
		// Successfully parsed as time, return as timeStamp
		ts := &timeStamp{start: t1, end: t2}
		return nil, ts, nil
	} else {
		// Successfully parsed as date
		return t, nil, nil
	}
}

// parseDate processes date expressions and returns a time.Time
// This function handles various date formats and keywords
//
// Example inputs:
//  1. Keywords:
//     - "eod" (end of day) -> today at 23:59
//     - "now" -> current time
//     - "tmrw" -> tomorrow
//     - "yest" -> yesterday
//  2. Days of week:
//     - "mon", "tue", "wed", "thu", "fri"
//     - "sat" or "eow" (end of week)
//     - "sun"
//  3. Custom date formats (defined in dateFormats)
func parseDate(arg string) (*time.Time, error) {
	// Get current time and end of today (23:59)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.UTC)

	// Process different date keywords
	switch arg {
	case "eod": // End of Day
		return &today, nil
	case "now": // Current time
		return &now, nil
	case "tmrw": // Tomorrow
		d := today.AddDate(0, 0, 1)
		return &d, nil
	case "yest": // Yesterday
		d := today.AddDate(0, 0, -1)
		return &d, nil
	case "eow", "sat": // End of Week (Saturday)
		// Calculate days until next Saturday
		// Example: If today is Wednesday (3), then Saturday (6) - Wednesday (3) = 3 days
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
		// If not a keyword, try parsing with predefined date formats
		for _, v := range dateFormats {
			t, err := time.Parse(v, arg)
			if err == nil {
				return &t, nil
			}
		}
		return nil, fmt.Errorf("invalid date format: %s", arg)
	}
}

// parseTimeStamp processes time expressions in 12-hour format
// This function handles complex time formats including ranges
//
// Valid formats:
//  1. Single time:
//     - "2pm", "2:30pm", "11am"
//  2. Time ranges:
//     - "2pm-4pm"
//     - "2:30pm-4:30pm"
//     - "2-4pm" (implicitly means 2pm-4pm)
//
// Rules:
//   - Hours must be 1-12
//   - Minutes must be 00-59
//   - AM/PM is case-insensitive
//   - Minutes must have 2 digits
//   - Colons must come between hours and minutes
func parseTimeStamp(arg string) (*time.Time, *time.Time, error) {
	// Flags to track what we've seen
	hour, colon, minute, am, dash := false, false, false, false, false
	// Values to store parsed time components
	startHour, endHour, startMinute, endMinute := -1, -1, -1, -1
	// Store AM/PM format for start and end times
	startFormat, endFormat := "", ""

	// Parse the time string character by character
	// Example: "2:30pm-4:45pm"
	for i := 0; i < len(arg); i++ {
		switch arg[i] {
		case ':':
			// Handle colon between hours and minutes (e.g., "2:30")
			if !hour || minute {
				return nil, nil, fmt.Errorf("colon is expected after hour, not minute: %s", arg)
			} else if am {
				return nil, nil, fmt.Errorf("colon can never occur after 'am' or 'pm': %s", arg)
			} else if colon {
				return nil, nil, fmt.Errorf("colons cannot be duplicated for same hour, minute combination: %s", arg)
			}
			colon = true

		case '-':
			// Handle dash between start and end times (e.g., "2-4")
			if !hour || dash {
				return nil, nil, fmt.Errorf("dashes require an hour and cannot be duplicated: %s", arg)
			}
			// Reset flags for parsing end time
			dash = true
			hour = false
			colon = false
			minute = false
			am = false

		case 'a', 'A', 'p', 'P':
			// Handle AM/PM indicators
			if !hour || am {
				return nil, nil, fmt.Errorf("'am'/'pm' cannot occur without an hour nor can there be duplicates: %s", arg)
			}
			// Check for 'm' or 'M' after 'a' or 'p'
			if i+1 < len(arg) && (arg[i+1] == 'm' || arg[i+1] == 'M') {
				s := arg[i : i+2] // Get "am" or "pm"
				if !dash {
					startFormat = s // Store format for start time
				} else {
					endFormat = s // Store format for end time
				}
			} else {
				return nil, nil, fmt.Errorf("provided time tag besides valid 'am'/'pm': %s", arg)
			}
			am = true
			i++ // Skip the 'm' we just processed

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// Handle numeric components (hours and minutes)
			if hour && !colon {
				// If we've seen an hour but no colon, we're in an invalid state
				return nil, nil, fmt.Errorf("minutes must be separated by colon from hours: %s", arg)
			} else if minute {
				// Can't have more than one set of minutes
				return nil, nil, fmt.Errorf("minutes cannot be duplicated: %s", arg)
			} else if am {
				// Numbers can't come after AM/PM
				return nil, nil, fmt.Errorf("hours and minutes cannot come after 'am'/'pm' signature: %s", arg)
			} else if hour {
				// Processing minutes (must be two digits)
				minute = true
				if i+1 < len(arg) && arg[i+1] >= '0' && arg[i+1] <= '9' {
					// Parse two-digit minutes
					x, _ := strconv.Atoi(arg[i : i+2])
					if x > 59 {
						return nil, nil, fmt.Errorf("minutes cannot be greater than 59: %s", arg)
					}
					// Store minutes for start or end time
					if !dash {
						startMinute = x
					} else {
						endMinute = x
					}
					i++ // Skip the second digit we just processed
				} else {
					return nil, nil, fmt.Errorf("minutes require 2 digits: %s", arg)
				}
			} else {
				// Processing hours (can be one or two digits)
				hour = true
				if i+1 >= len(arg) || arg[i+1] < '0' || arg[i+1] > '9' {
					// Single-digit hour
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
					// Two-digit hour
					x, _ := strconv.Atoi(arg[i : i+2])
					if x > 12 {
						return nil, nil, fmt.Errorf("hours cannot be greater than 12: %s", arg)
					}
					if !dash {
						startHour = x
					} else {
						endHour = x
					}
					i++ // Skip the second digit we just processed
				} else {
					return nil, nil, fmt.Errorf("minutes require 2 digits: %s", arg)
				}
			}
		default:
			// Any other character is invalid
			return nil, nil, fmt.Errorf("invalid time format: %s", arg)
		}
	}

	// Convert times from 12-hour to 24-hour format

	// Handle start time AM/PM conversion
	if (startFormat == "pm" || startFormat == "PM") && startHour != 12 {
		startHour += 12 // 2pm becomes 14
	} else if (startFormat == "am" || startFormat == "AM") && startHour == 12 {
		startHour = 0 // 12am becomes 00
	}

	// Set default minutes if not specified
	if startMinute == -1 {
		startMinute = 0
	}

	// If we only have a start time (no range)
	if endHour == -1 {
		// Create time object for the start time only
		t := time.Now()
		st := time.Date(t.Year(), t.Month(), t.Day(), startHour, startMinute, 0, 0, time.UTC)
		return &st, nil, nil
	}

	// Handle end time AM/PM conversion
	if (endFormat == "pm" || endFormat == "PM") && endHour != 12 {
		endHour += 12 // 4pm becomes 16
	} else if (endFormat == "am" || endFormat == "AM") && endHour == 12 {
		endHour = 0 // 12am becomes 00
	} else if endFormat == "" && endHour < startHour {
		// If no AM/PM specified for end time and it's less than start,
		// assume PM (e.g., "2-4" means "2pm-4pm")
		endHour += 12
	}

	// Set default end minutes if not specified
	if endMinute == -1 {
		endMinute = 0
	}

	// Validate that end time is after start time
	if startHour > endHour || (startHour == endHour && startMinute >= endMinute) {
		return nil, nil, fmt.Errorf("starting time must be earlier than ending time: %s", arg)
	}

	// Create time objects for both start and end times
	t := time.Now()
	st := time.Date(t.Year(), t.Month(), t.Day(), startHour, startMinute, 0, 0, time.UTC)
	et := time.Date(t.Year(), t.Month(), t.Day(), endHour, endMinute, 0, 0, time.UTC)
	return &st, &et, nil
}
