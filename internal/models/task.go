package models

import (
	"fmt"
	"time"
)

type Task struct {
	ID       int
	Desc     string
	TS       *TimeStamp
	Tags     []string
	Priority int
}

type TimeStamp struct {
	Start *time.Time
	End   *time.Time
}

func (t *TimeStamp) String() string {
	return fmt.Sprintf("%s: %s-%s", t.Start.Format("02-01-2006"), t.Start.Format(time.Kitchen), t.End.Format(time.Kitchen))
}

func NewTask(id int, description string, ts *TimeStamp, tags []string, priority int) *Task {
	return &Task{ID: id, Desc: description, TS: ts, Tags: tags, Priority: priority}
}
