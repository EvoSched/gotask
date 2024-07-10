package repository

import (
	"database/sql"
	"time"

	"github.com/EvoSched/gotask/internal/models"
)

var due time.Time = time.Now()

var tasks = []*models.Task{
	models.NewTask(1, "title1", "description1", &due, []string{"MA"}),
	models.NewTask(2, "title2", "description2", nil, []string{"CS"}),
	models.NewTask(3, "title3", "description3", &due, []string{"MA"}),
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
