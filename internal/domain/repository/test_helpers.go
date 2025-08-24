package repository

import (
	"database/sql"
	"github.com/pressly/goose"
)

func applyMigrations(migrationPath string, db *sql.DB) (err error) {
	if err = goose.SetDialect("postgres"); err != nil {
		return
	}

	if err = goose.Up(db, migrationPath); err != nil {
		return
	}

	return
}
