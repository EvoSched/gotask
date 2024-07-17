package models

import (
	"fmt"
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
