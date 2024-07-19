package models

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID        int
	Desc      string
	Priority  int
	Tags      []string   // tags, tags: string tag1,tag2,tag3
	Notes     []string   // comment1,comment2,comment3
	StartAt   *time.Time // timestamp datetime
	EndAt     *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Finished  bool
}

//tasks
//tags
// tag1

//task_tags
//tag_id, taks_id

func NewTask(id int, desc string, priority int, tags []string, comments []string, startAt *time.Time, endAt *time.Time) *Task {
	now := time.Now()
	return &Task{
		ID:        id,
		Desc:      desc,
		Priority:  priority,
		Tags:      tags,
		Notes:     comments,
		StartAt:   startAt,
		EndAt:     endAt,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func DisplayTask(task *Task) {
	// Print header
	fmt.Println("Task Details:")
	fmt.Println("--------------")
	fmt.Printf("ID             %d\n", task.ID)
	fmt.Printf("Description    %s\n", task.Desc)
	fmt.Printf("Priority       %d\n", task.Priority)
	if len(task.Tags) > 0 {
		fmt.Printf("Tags           %v\n", task.Tags)
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

	fmt.Printf("\nNotes:\nThu, 18 Jul 2024 00:50:46 EDT - unexpected issue came up, am resolving now")
}

// / FormatTask prints a task in the desired format
func FormatTask(task *Task) string {
	// Format the status
	status := "[ ]"
	if task.Finished {
		status = "[x]"
	}

	// Format the tags
	tags := strings.Join(task.Tags, ", ")

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
	return fmt.Sprintf("%d   %s     %-30s %d          %-13s %s   ", // Adjusted format string with extra spaces
		task.ID, status, task.Desc, task.Priority, tags, due)
}

// DisplayTasks prints a list of tasks in the desired format
func DisplayTasks(tasks []*Task) {
	// Print header
	fmt.Println("ID  Status  Desc                           Priority   Tags          Due   ")
	fmt.Println("------------------------------------------------------------------------------------------------------")

	// Print each task
	for _, task := range tasks {
		fmt.Println(FormatTask(task))
	}
}
