package service

import (
	"database/sql"
	"github.com/EvoSched/gotask/internal/sqlite"
	"github.com/EvoSched/gotask/internal/types"
)

type TaskRepoQuery interface {
	GetTask(id int) (*types.Task, error)
	GetTasks() ([]*types.Task, error)
	GetDesc(id int) (string, error)
	GetNotes(id int) ([]string, error)
}

type TaskRepoStmt interface {
	AddTask(task *types.Task) (int, error)
	AddNote(id int, note string) error
	UpdateStatus(id int, status bool) error
	UpdateTask(task *types.Task) error
	DeleteTask(db *sql.DB, id int) error
}

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db}
}

func (r *TaskRepo) GetDesc(id int) (string, error) {
	return sqlite.QueryDesc(r.db, id)
}

func (r *TaskRepo) GetNotes(id int) ([]string, error) {
	return sqlite.QueryNotes(r.db, id)
}

func (r *TaskRepo) GetTask(id int) (*types.Task, error) {
	t, err := sqlite.QueryTask(r.db, id)
	if err != nil {
		return nil, err
	}
	n, err := sqlite.QueryNotes(r.db, id)
	if err != nil {
		return nil, err
	}
	t.Notes = append(t.Notes, n...)
	return &t, nil
}

func (r *TaskRepo) GetTasks() ([]*types.Task, error) {
	return sqlite.QueryTasks(r.db)
}

func (r *TaskRepo) AddTask(task *types.Task) (int, error) {
	err := sqlite.InsertTask(r.db, task)
	if err != nil {
		return 0, err
	}
	return sqlite.QueryLastID(r.db)
}

func (r *TaskRepo) AddNote(id int, note string) error {
	return sqlite.InsertNote(r.db, id, note)
}

func (r *TaskRepo) UpdateStatus(id int, status bool) error {
	return sqlite.UpdateStatus(r.db, id, status)
}

func (r *TaskRepo) UpdateTask(task *types.Task) error {
	return sqlite.UpdateTask(r.db, task)
}

func (r *TaskRepo) DeleteTask(id int) error {
	return sqlite.DeleteTask(r.db, id)
}
