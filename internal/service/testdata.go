package service

import (
	"github.com/EvoSched/gotask/internal/types"
	"time"
)

var curr = time.Now()
var start = time.Date(curr.Year(), curr.Month(), curr.Day(), 13, 15, 0, 0, time.UTC)
var end = time.Date(curr.Year(), curr.Month(), curr.Day(), 15, 30, 0, 0, time.UTC)
var date = time.Date(curr.Year(), curr.Month(), curr.Day(), 23, 59, 0, 0, time.UTC)

// sample data to test command functions
var tasks = []*types.Task{
	types.NewTask(1, "finish project3", 5, []string{"MA", "CS"}, []string{"comment1"}, &start, nil),
	types.NewTask(2, "study BSTs", 8, []string{"CS"}, []string{"comment2"}, &start, &end),
	types.NewTask(3, "lunch with Edgar", 2, []string{"Fun"}, []string{"comment3"}, nil, nil),
	types.NewTask(4, "meeting for db proposal", 5, []string{"Project"}, []string{"comment4"}, &date, nil),
}
