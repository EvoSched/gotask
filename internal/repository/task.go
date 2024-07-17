package repository

import (
	"database/sql"

	"github.com/EvoSched/gotask/internal/models"
)

// sample data to test command functions
var tasks = []*models.Task{
	models.NewTask(1, "description1", 5, []string{"MA", "CS"}, []string{"comment1"}, nil, nil),
	models.NewTask(2, "description2", 8, []string{"CS"}, []string{"comment2"}, nil, nil),
	models.NewTask(3, "description3", 2, []string{"MA", "CS"}, []string{"comment3"}, nil, nil),
	models.NewTask(4, "description4", 5, []string{"CH"}, []string{"comment4"}, nil, nil),
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
