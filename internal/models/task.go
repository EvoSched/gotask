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

func DisplayTasks(task []*Task) {
	// Print header
	fmt.Printf("%-5s %-30s %-10s %-15s %-25s %-25s\n", "ID", "Desc", "Priority", "Tags", "StartAt", "EndAt")
	for _, t := range task {
		fmtTask(t)
	}
}

func fmtTask(task *Task) {
	// Print task details
	tags := strings.Join(task.Tags, ", ")
	fmt.Printf("%-5d %-30s %-10d %-15s", task.ID, task.Desc, task.Priority, tags)
	if task.StartAt != nil {
		fmt.Printf(" %-25s", task.StartAt.Format(time.RFC3339))
	} else {
		fmt.Printf(" %-25s", "N/A")
	}
	if task.EndAt != nil {
		fmt.Printf(" %-25s", task.EndAt.Format(time.RFC3339))
	} else {
		fmt.Printf(" %-25s", "N/A")
	}
	fmt.Println()
}
