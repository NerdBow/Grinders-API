package sqlite

import (
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddUser(user util.User) error {
	query := "INSERT INTO users (username, hash, creation_time) VALUES (?, ?, ?);"
	result, err := db.Exec(query, user.Username, user.Hash, user.CreationTime)
	if err != nil {
		slog.Error("")
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		slog.Warn("")
	}
	if n < 1 {
		slog.Warn("")
	}

	return nil
}

func (db *SQLiteDB) GetUser(userId uint64) (util.User, error) {
	query := "SELECT id, username, hash, creation_time FROM users WHERE id = ?;"
	row := db.QueryRow(query, userId)

	user := util.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Hash, &user.CreationTime)
	if err != nil {
		slog.Error("")
		return user, err
	}
	return user, nil
}

func (db *SQLiteDB) GetUserByUsername(username string) (util.User, error) {
	query := "SELECT id, username, hash, creation_time FROM users WHERE username = ?;"
	row := db.QueryRow(query, username)

	user := util.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Hash, &user.CreationTime)
	if err != nil {
		slog.Error("")
		return user, err
	}
	return user, nil
}

func (db *SQLiteDB) EditUsername(userId uint64, newName string) error {
	query := "UPDATE users SET username = ? WHERE id = ?;"
	result, err := db.Exec(query, newName, userId)
	if err != nil {
		slog.Error("")
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		slog.Warn("")
	}
	if n < 1 {
		slog.Warn("")
	}

	return nil
}
