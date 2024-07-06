package models

type Event struct {
	ID          int
	Title       string
	Description string
}

func NewEvent(id int, title, description string) *Event {
	return &Event{ID: id, Title: title, Description: description}
}
