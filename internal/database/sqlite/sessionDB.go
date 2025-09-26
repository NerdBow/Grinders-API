package sqlite

import (
	"context"
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddSession(logger *slog.Logger, session util.Session) error {
	query := "INSERT INTO sessions (id, expiration_time, creation_time, user_id) VALUES (?, ?, ?, ?);"
	result, err := db.Exec(query, session.HashedId, session.ExpirationTime, session.CreationTime, session.UserId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec AddSession", slog.String("err", err.Error()))
		return util.ErrDatabase
	}
	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddSession", slog.String("err", err.Error()))
	}
	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddSession", slog.String("err", "There were no rows affected"))
	}
	return nil
}

func (db *SQLiteDB) GetSession(logger *slog.Logger, hashedId string, userId uint64) (util.Session, error) {
	query := "SELECT id, expiration_time, creation_time, user_id FROM sessions WHERE id = ? AND user_id = ?;"
	rows, err := db.Query(query, hashedId, userId)
	session := util.Session{}

	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Query GetSession", slog.String("err", err.Error()))
		return session, util.ErrDatabase
	}

	for rows.Next() {
		err = rows.Scan(&session.HashedId, &session.ExpirationTime, &session.CreationTime, &session.UserId)
		if err != nil {
			logger.LogAttrs(context.Background(), slog.LevelError, "Scan GetSession", slog.String("err", err.Error()))
			return session, util.ErrDatabase
		}
	}
	return session, nil
}

func (db *SQLiteDB) DeleteSession(logger *slog.Logger, hashedId string) error {
	query := "DELETE FROM sessions WHERE id = ?"
	result, err := db.Exec(query, hashedId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec DeleteSession", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected DeleteSession", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected DeleteSession", slog.String("err", "There were no rows affected"))
	}

	return nil
}
