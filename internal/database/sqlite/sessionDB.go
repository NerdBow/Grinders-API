package sqlite

import (
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddSession(session util.Session) error {
	query := "INSERT INTO sessions (id, expiration_time, creation_time, user_id) VALUES (?, ?, ?, ?);"
	result, err := db.Exec(query, session.HashedId, session.ExpirationTime, session.CreationTime, session.UserId)
	if err != nil {
		slog.Error("")
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		slog.Error("")
		return nil
	}
	if n != 1 {
		slog.Error("")
	}
	return nil
}

func (db *SQLiteDB) GetSession(hashedId string, userId uint64) (util.Session, error) {
	query := "SELECT id, expiration_time, creation_time, user_id FROM sessions WHERE id = ? AND user_id = ?;"
	rows, err := db.Query(query, hashedId, userId)
	session := util.Session{}

	if err != nil {
		slog.Error("Query Error")
		return session, err
	}

	for rows.Next() {
		err = rows.Scan(&session.HashedId, &session.ExpirationTime, &session.CreationTime, &session.UserId)
		if err != nil {
			slog.Error("Scan Error")
			return session, err
		}
	}
	return session, nil
}

func (db *SQLiteDB) DeleteSession(hashedId string) error {
	query := "DELETE FROM sessions WHERE id = ?"
	result, err := db.Exec(query, hashedId)
	if err != nil {
		slog.Error("")
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		slog.Error("")
		return nil
	}

	if n != 1 {
		slog.Error("")
	}

	return nil
}
