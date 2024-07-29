package sqlite

import (
	"database/sql"
	"github.com/EvoSched/gotask/internal/config"
	"github.com/EvoSched/gotask/internal/types"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"os"
)

const (
	SQLiteDriver = "sqlite3"
)

func NewSQLite(config *config.SQLite) (*sql.DB, error) {
	// todo create database if it doesn't already exist
	err := setupDB(config)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(SQLiteDriver, config.Database)
	if err != nil {
		return nil, err
	}
	// we just created our database file, need to create tables now
	err = createTaskTable(db)
	if err != nil {
		return nil, err
	}
	err = createNoteTable(db)
	if err != nil {
		return nil, err
	}
	return db, err
}

func setupDB(config *config.SQLite) error {
	if _, err := os.Stat(config.Database); os.IsNotExist(err) {
		file, err := os.Create(config.Database)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func createTaskTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS task (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
	"desc" TEXT NOT NULL,
	"priority" INTEGER NOT NULL,
	"start_at" DATETIME,
	"end_at" DATETIME,
	"updated_at" DATETIME NOT NULL,
	"finished" INTEGER NOT NULL CHECK (finished IN (0,1))
);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func createNoteTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS note (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "task_id" INTEGER NOT NULL,
    "comment" TEXT NOT NULL
);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func QueryTask(db *sql.DB, id int) (types.Task, error) {
	var task types.Task
	row := db.QueryRow(`SELECT id, desc, priority, start_at, end_at, updated_at, finished FROM task WHERE id = ?`, id)
	err := row.Scan(&task.ID, &task.Desc, &task.Priority, &task.StartAt, &task.EndAt, &task.UpdatedAt, &task.Finished)
	if err != nil {
		return task, err
	}
	return task, nil
}

func QueryTasks(db *sql.DB) ([]*types.Task, error) {
	rows, err := db.Query(`SELECT id, desc, priority, start_at, end_at, updated_at, finished FROM task`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*types.Task
	for rows.Next() {
		var task types.Task
		err := rows.Scan(&task.ID, &task.Desc, &task.Priority, &task.StartAt, &task.EndAt, &task.UpdatedAt, &task.Finished)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func QueryDesc(db *sql.DB, id int) (string, error) {
	var desc string
	row := db.QueryRow(`SELECT desc FROM task WHERE id = ?`, id)
	err := row.Scan(&desc)
	if err != nil {
		return "", err
	}
	return desc, nil
}

func QueryNotes(db *sql.DB, id int) ([]string, error) {
	rows, err := db.Query(`SELECT comment FROM note WHERE task_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var n string
		err := rows.Scan(&n)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func QueryLastID(db *sql.DB) (int, error) {
	// Query to find the maximum existing ID
	row := db.QueryRow(`SELECT COALESCE(MAX(id), 0) FROM task`)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func InsertTask(db *sql.DB, task *types.Task) error {
	stmt, err := db.Prepare(`INSERT INTO task(desc, priority, start_at, end_at, updated_at, finished) VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Desc, task.Priority, task.StartAt, task.EndAt, task.UpdatedAt, task.Finished)
	return err
}

func InsertNote(db *sql.DB, id int, note string) error {
	stmt, err := db.Prepare(`INSERT INTO note(task_id, comment) VALUES(?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, note)
	return err
}

func UpdateStatus(db *sql.DB, id int, status bool) error {
	stmt, err := db.Prepare(`UPDATE task SET finished = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, id)
	return err
}

func UpdateTask(db *sql.DB, task *types.Task) error {
	stmt, err := db.Prepare(`UPDATE task SET desc = ?, priority = ?, start_at = ?, end_at = ?, updated_at = ?, finished = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(task.Desc, task.Priority, task.StartAt, task.EndAt, task.UpdatedAt, task.Finished, task.ID)
	return err
}

func DeleteTask(db *sql.DB, id int) error {
	stmt, err := db.Prepare(`DELETE FROM task WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}
