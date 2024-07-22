package repository

import (
	"database/sql"

	"github.com/EvoSched/gotask/internal/models"
)

type Task interface {
	GetTask(id int) (*models.Task, error)
	GetTasks() ([]*models.Task, error)
}

type Repository struct {
	Task
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{Task: NewTaskRepository(db)}
}
