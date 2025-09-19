package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	*sql.DB
}

func NewSQLiteDB(file string) (SQLiteDB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal=WAL", file))
	if err != nil {
		return SQLiteDB{}, err
	}
	return SQLiteDB{db}, nil
}

func (db *SQLiteDB) CreateTables() error {
	query := `
CREATE TABLE IF NOT EXISTS "sessions" (
	"id" TEXT NOT NULL UNIQUE,
	"expiration_time" TIMESTAMP NOT NULL,
	"creation_time" TIMESTAMP NOT NULL,
	"user_id" INTEGER NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY ("user_id") REFERENCES "users"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);
`
	_, err := db.Exec(query)
	if err != nil {
		slog.Error("")
		return err
	}
	return nil
}
