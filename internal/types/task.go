package types

import (
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
	UpdatedAt *time.Time
	Finished  bool
}

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
		UpdatedAt: &now,
	}
}
