package repository

import (
	"database/sql"

	"github.com/EvoSched/gotask/internal/config"
)

const (
	SQLiteDriver = "sqlite3"
)

func NewSQLite(config *config.SQLite) (*sql.DB, error) {
	return sql.Open(SQLiteDriver, config.Database)
}
