package sqlite

import (
	"context"
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddUser(logger *slog.Logger, user util.User) error {
	query := "INSERT INTO users (username, hash, creation_time) VALUES (?, ?, ?);"
	result, err := db.Exec(query, user.Username, user.Hash, user.CreationTime)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec AddUser", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddUser", slog.String("err", err.Error()))
	}
	if n < 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddUser", slog.String("err", "There were no rows affected"))
	}

	return nil
}

func (db *SQLiteDB) GetUser(logger *slog.Logger, userId uint64) (util.User, error) {
	query := "SELECT id, username, hash, creation_time FROM users WHERE id = ?;"
	row := db.QueryRow(query, userId)

	user := util.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Hash, &user.CreationTime)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Scan GetUser", slog.String("err", err.Error()))
		return user, util.ErrDatabase
	}
	return user, nil
}

func (db *SQLiteDB) GetUserByUsername(logger *slog.Logger, username string) (util.User, error) {
	query := "SELECT id, username, hash, creation_time FROM users WHERE username = ?;"
	row := db.QueryRow(query, username)

	user := util.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Hash, &user.CreationTime)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Scan GetUserByUsername", slog.String("err", err.Error()))
		return user, util.ErrDatabase
	}
	return user, nil
}

func (db *SQLiteDB) EditUsername(logger slog.Logger, userId uint64, newName string) error {
	query := "UPDATE users SET username = ? WHERE id = ?;"
	result, err := db.Exec(query, newName, userId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec EditUsername", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected EditUsername", slog.String("err", err.Error()))
	}
	if n < 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected EditUsername", slog.String("err", "There were no rows affected"))
	}

	return nil
}
