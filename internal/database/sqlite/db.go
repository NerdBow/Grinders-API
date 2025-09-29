package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	*sql.DB
}

func NewSQLiteDB(file string) (SQLiteDB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal=WAL", file))
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "SQLiteDB Open", slog.String("err", err.Error()))
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
CREATE TABLE IF NOT EXISTS "users" (
	"id" INTEGER NOT NULL UNIQUE,
	"username" TEXT NOT NULL UNIQUE,
	"hash" TEXT NOT NULL,
	"creation_time" TIMESTAMP,
	PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "categories" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" TEXT NOT NULL,
	"user_id" INTEGER NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY ("user_id") REFERENCES "users"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE TABLE IF NOT EXISTS "tasks" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" TEXT NOT NULL,
	"creation_time" TIMESTAMP NOT NULL,
	"completion_time" TIMESTAMP NOT NULL,
	"deadline_time" TIMESTAMP,
	"is_completed" BOOLEAN NOT NULL,
	"category_id" INTEGER NOT NULL,
	"user_id" INTEGER NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY ("user_id") REFERENCES "users"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("category_id") REFERENCES "category"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);
`
	_, err := db.Exec(query)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "SQLiteDB Create Table", slog.String("err", err.Error()))
		return util.ErrDatabase
	}
	return nil
}
