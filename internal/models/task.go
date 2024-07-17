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
	Comments  []string   // comment1,comment2,comment3
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
		Comments:  comments,
		StartAt:   startAt,
		EndAt:     endAt,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func DisplayTask(task *Task) {
	// Print header
	fmt.Printf("%-5s %-30s %-10s %-15s %-25s %-25s\n", "ID", "Desc", "Priority", "Tags", "StartAt", "EndAt")
	fmtTask(task)
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
