package repository

import (
	"database/sql"

	"github.com/EvoSched/gotask/internal/models"
)

var tasks = []*models.Task{
	models.NewTask(1, "title1", "description1"),
	models.NewTask(2, "title2", "description2"),
	models.NewTask(3, "title3", "description3"),
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
