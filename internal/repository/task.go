package repository

import (
	"database/sql"
	"time"

	"github.com/EvoSched/gotask/internal/models"
)

var curr = time.Now()
var start = time.Date(curr.Year(), curr.Month(), curr.Day(), 13, 15, 0, 0, time.UTC)
var end = time.Date(curr.Year(), curr.Month(), curr.Day(), 15, 30, 0, 0, time.UTC)
var date = time.Date(curr.Year(), curr.Month(), curr.Day(), 23, 59, 0, 0, time.UTC)

// sample data to test command functions
var tasks = []*models.Task{
	models.NewTask(1, "finish project3", 5, []string{"MA", "CS"}, []string{"comment1"}, &start, nil),
	models.NewTask(2, "study BSTs", 8, []string{"CS"}, []string{"comment2"}, &start, &end),
	models.NewTask(3, "lunch with Edgar", 2, []string{"Fun"}, []string{"comment3"}, nil, nil),
	models.NewTask(4, "meeting for db proposal", 5, []string{"Project"}, []string{"comment4"}, &date, nil),
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) GetTask(id int) (*models.Task, error) {
	//TODO: sql query

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, nil
}

func (r *TaskRepository) GetTasks() ([]*models.Task, error) {
	//TODO: sql query

	return tasks, nil
}
