package models

type Task struct {
	ID          int
	Title       string
	Description string
}

func NewTask(id int, title, description string) *Task {
	return &Task{ID: id, Title: title, Description: description}
}
