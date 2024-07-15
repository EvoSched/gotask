package models

import "time"

type Task struct {
	ID          int
	Description string
	Date        *time.Time
	Tags        []string
}

func NewTask(id int, description string, date *time.Time, tags []string) *Task {
	return &Task{ID: id, Description: description, Date: date, Tags: tags}
}
