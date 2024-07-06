package service

import (
	"github.com/EvoSched/gotask/internal/models"
	"github.com/EvoSched/gotask/internal/repository"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo}
}

func (s *TaskService) GetTask(id int) (*models.Task, error) {
	return s.repo.GetTask(id)
}

func (s *TaskService) GetTasks() ([]*models.Task, error) {
	return s.repo.GetTasks()
}
