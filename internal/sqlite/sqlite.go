package sqlite

import (
	"database/sql"
	"os"

	"github.com/EvoSched/gotask/internal/config"
)

const (
	SQLiteDriver = "sqlite3"
)

func NewSQLite(config *config.SQLite) (*sql.DB, error) {
	// todo create database if it doesn't already exist
	flag, err := setupDB(config)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(SQLiteDriver, config.Database)
	if err != nil {
		return nil, err
	}
	// we just created our database file, need to create tables now
	if flag {
		err = createTable(db)
	}
	return db, err
}

func setupDB(config *config.SQLite) (bool, error) {
	if _, err := os.Stat(config.Database); os.IsNotExist(err) {
		file, err := os.Create(config.Database)
		if err != nil {
			return false, err
		}
		file.Close()
		return true, nil
	}
	return false, nil
}

func createTable(db *sql.DB) error {
	tableSQL := `CREATE TABLE task (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"desc" TEXT NOT NULL,
		"priority" INTEGER NOT NULL,
		"start_at" DATETIME,
		"end_at" DATETIME,
		"created_at" DATETIME NOT NULL ,
		"finished" INTEGER NOT NULL CHECK (finished IN (0,1))
	  );`

	statement, err := db.Prepare(tableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}

	// Execute SQL Statements
	_, err = statement.Exec()
	return err
}
