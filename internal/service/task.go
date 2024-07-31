package service

import (
	"database/sql"
	"github.com/EvoSched/gotask/internal/sqlite"
	"github.com/EvoSched/gotask/internal/types"
	"time"
)

type TaskRepoQuery interface {
	GetTask(id int) (*types.Task, error)
	GetTasks() ([]*types.Task, error)
	GetDesc(id int) (string, error)
	GetTasksDue() ([]*types.Task, error)
	GetTasksArchived() ([]*types.Task, error)
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
	return sqlite.QueryTaskDesc(r.db, id)
}

func (r *TaskRepo) GetTask(id int) (*types.Task, error) {
	t, err := sqlite.QueryTask(r.db, id)
	if err != nil {
		return nil, err
	}
	n, err := sqlite.QueryTaskNotes(r.db, id)
	if err != nil {
		return nil, err
	}
	t.Notes = append(t.Notes, n...)
	tags, err := sqlite.QueryTaskTags(r.db, id)
	if err != nil {
		return nil, err
	}
	t.Tags = append(t.Tags, tags...)
	return &t, nil
}

func (r *TaskRepo) GetTasks() ([]*types.Task, error) {
	tasks, err := sqlite.QueryTasks(r.db)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		tags, err := sqlite.QueryTaskTags(r.db, t.ID)
		if err != nil {
			return nil, err
		}
		t.Tags = append(t.Tags, tags...)
	}
	return tasks, nil
}

func (r *TaskRepo) GetTasksDue() ([]*types.Task, error) {
	tasks, err := sqlite.QueryTasksArchived(r.db, false)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		tags, err := sqlite.QueryTaskTags(r.db, t.ID)
		if err != nil {
			return nil, err
		}
		t.Tags = append(t.Tags, tags...)
	}
	return tasks, nil
}

func (r *TaskRepo) GetTasksArchived() ([]*types.Task, error) {
	tasks, err := sqlite.QueryTasksArchived(r.db, true)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		tags, err := sqlite.QueryTaskTags(r.db, t.ID)
		if err != nil {
			return nil, err
		}
		t.Tags = append(t.Tags, tags...)
	}
	return tasks, nil
}

func (r *TaskRepo) AddTask(task *types.Task) (int, error) {
	err := sqlite.InsertTask(r.db, task)
	if err != nil {
		return 0, err
	}
	i, err := sqlite.QueryLastID(r.db)
	if err != nil {
		return 0, err
	}
	for _, t := range task.Tags {
		ti, err := sqlite.QueryTag(r.db, t)
		if err != nil {
			err = sqlite.InsertTag(r.db, t)
			if err != nil {
				return 0, err
			}
			ti, err = sqlite.QueryTag(r.db, t)
			if err != nil {
				return 0, err
			}
		}
		err = sqlite.InsertTagPair(r.db, i, ti)
		if err != nil {
			return 0, err
		}
	}
	return sqlite.QueryLastID(r.db)
}

func (r *TaskRepo) AddNote(id int, note string) error {
	return sqlite.InsertNote(r.db, id, note)
}

func (r *TaskRepo) UpdateStatus(id int, status bool) error {
	err := sqlite.UpdateStatus(r.db, id, status)
	if err != nil {
		return err
	}
	t := time.Now()
	if status {
		err = sqlite.UpdateCompletedAt(r.db, id, &t)
	} else {
		err = sqlite.UpdateCompletedAt(r.db, id, nil)
	}
	return err
}

func (r *TaskRepo) UpdateTask(task *types.Task) error {
	err := sqlite.UpdateTask(r.db, task)
	if err != nil {
		return err
	}
	for _, t := range task.Tags {
		ti, err := sqlite.QueryTag(r.db, t)
		if err != nil {
			err = sqlite.InsertTag(r.db, t)
			if err != nil {
				return err
			}
			ti, err = sqlite.QueryTag(r.db, t)
			if err != nil {
				return err
			}
		}
		// check whether tag pair exists
		err = sqlite.QueryTagPair(r.db, task.ID, ti)
		if err != nil {
			err = sqlite.InsertTagPair(r.db, task.ID, ti)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *TaskRepo) DeleteTask(id int) error {
	return sqlite.DeleteTask(r.db, id)
}
