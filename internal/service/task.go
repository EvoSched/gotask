package service

import (
	"database/sql"
	"github.com/EvoSched/gotask/internal/types"
)

type TaskRepoQuery interface {
	GetTask(id int) (*types.Task, error)
	GetTasks() ([]*types.Task, error)
	// todo need to add interface functions for getting tasks by priority, tags, and time
}

type TaskRepoStmt interface {
	AddTask(task *types.Task) error
	EditTask(id int, task *types.Task) error
	RemoveTask(id int) error
	CompleteTask(id int) error
	IncompleteTask(id int) error
	AddNote(id int, note string) error
}

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db}
}

func (r *TaskRepo) GetTask(id int) (*types.Task, error) {
	//TODO: sql query

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, nil
}

func (r *TaskRepo) GetTasks() ([]*types.Task, error) {
	//TODO: sql query

	// todo remove when integrating SQL (this is purely for displaying mock data
	tasks[1].Finished = true
	tasks[3].Finished = true

	return tasks, nil
}
