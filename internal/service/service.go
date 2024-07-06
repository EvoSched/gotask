package service

import (
	"github.com/EvoSched/gotask/internal/models"
	"github.com/EvoSched/gotask/internal/repository"
)

type Task interface {
	GetTask(id int) (*models.Task, error)
	GetTasks() ([]*models.Task, error)
}

type Service struct {
	Task
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Task: NewTaskService(repo.Task),
	}
}
