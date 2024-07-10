package models

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Date        *time.Time
	Tags        []string
}

func NewTask(id int, title, description string, date *time.Time, tags []string) *Task {
	return &Task{ID: id, Title: title, Description: description, Date: date, Tags: tags}
}
