package sqlite

import (
	"context"
	"log/slog"
	"time"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddTask(logger *slog.Logger, task util.Task) error {
	query := `INSERT INTO tasks 
	(name, creation_time, deadline_time, completion_time, is_completed, category_id, user_id) VALUES 
	(?, ?, ?, ?, ?, ?, ?);`

	result, err := db.Exec(query, task.Name, task.CreationTime, task.DeadlineTime, time.Time{}, false, task.CategoryId, task.UserId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec AddTask", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddTask", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddTask", slog.String("err", "There were no rows affected"))
	}

	return nil
}

func (db *SQLiteDB) GetTask(logger *slog.Logger, taskId uint64, userId uint64) (util.Task, error) {
	query := "SELECT id, name, creation_time, completion_time, deadline_time, is_completed, category_id, user_id FROM tasks WHERE id = ? AND user_id = ?;"
	row := db.QueryRow(query, taskId, userId)

	task := util.Task{}
	err := row.Scan(&task.Id, &task.Name, &task.CreationTime, &task.CompletionTime, &task.DeadlineTime, &task.IsComplete, &task.CategoryId, &task.UserId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Scan GetTask", slog.String("err", err.Error()))
		return task, err
	}
	return task, nil
}

func (db *SQLiteDB) QueryTask(logger *slog.Logger, querySettings util.TaskQuerySettings) ([]util.Task, error) {
	query := `SELECT 
	id, name, creation_time, completion_time, deadline_time, is_completed, category_id, user_id
	FROM tasks 
	WHERE user_id = ?`

	params := make([]any, 0, 3)
	params = append(params, querySettings.UserId)

	if querySettings.Category != "" {
		query = `SELECT 
		t.id, t.name, t.creation_time, t.completion_time, t.deadline_time, t.is_completed, t.category_id, t.user_id
		FROM tasks t INNER JOIN categories c ON t.category_id == c.id
		WHERE t.user_id = ? AND c.name LIKE ?`
		params = append(params, "%"+querySettings.Category+"%")
	}

	if querySettings.Name != "" {
		query += " AND name LIKE ?"
		params = append(params, "%"+querySettings.Name+"%")
	}

	if querySettings.SortOrder != 0 && querySettings.SortType != 0 {
		query += " ORDER BY"

		switch querySettings.SortType {
		case util.SORT_COMPLETION:
			query += " completion_time"
		case util.SORT_CREATION:
			query += " creation_time"
		case util.SORT_DEADLINE:
			query += " deadline_time"
		}

		switch querySettings.SortOrder {
		case util.ORDER_ASCEDNING:
			query += " ASC"
		case util.ORDER_DESCEDNING:
			query += " DESC"
		}
	}
	query += " LIMIT ?,20;"

	params = append(params, (querySettings.Page-1)*PAGE_SIZE)

	logger.LogAttrs(context.Background(), slog.LevelDebug, "SQL Query QueryTask", slog.String("query", query))
	rows, err := db.Query(query, params...)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Query QueryTask", slog.String("err", err.Error()))
		return nil, err
	}

	tasks := make([]util.Task, 0, 20)
	for rows.Next() {
		task := util.Task{}

		err = rows.Scan(&task.Id, &task.Name, &task.CreationTime, &task.CompletionTime, &task.DeadlineTime, &task.IsComplete, &task.CategoryId, &task.UserId)
		if err != nil {
			logger.LogAttrs(context.Background(), slog.LevelError, "Scan QueryTask", slog.String("err", err.Error()))
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *SQLiteDB) EditTask(logger *slog.Logger, task util.Task) error {
	query := "UPDATE categories SET"

	params := make([]any, 0, 7)
	if task.Name != "" {
		query += " name = ?"
		params = append(params, task.Name)
	}
	if task.CreationTime.Equal(time.Time{}) {
		query += " creation_time = ?"
		params = append(params, task.CreationTime)
	}
	if task.CompletionTime.Equal(time.Time{}) {
		query += " completion_time = ?"
		params = append(params, task.CompletionTime)
	}
	if task.DeadlineTime.Equal(time.Time{}) {
		query += " deadline_time = ?"
		params = append(params, task.DeadlineTime)
	}
	if task.CategoryId != 0 {
		query += " category_id = ?"
		params = append(params, task.CategoryId)
	}
	query += " WHERE user_id = ? AND id = ?;"

	params = append(params, task.UserId, task.Id)

	result, err := db.Exec(query, params...)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec EditTask", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected EditTask", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected EditTask", slog.String("err", "There were no rows affected"))
	}

	return nil
}

func (db *SQLiteDB) DeleteTask(logger *slog.Logger, taskId uint64, userId uint64) error {
	query := "DELETE FROM tasks WHERE user_id = ? AND id = ?;"

	result, err := db.Exec(query, userId, taskId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec DeleteTask", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected DeleteTask", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected DeleteTask", slog.String("err", "There were no rows affected"))
	}

	return nil
}

func (db *SQLiteDB) SetTaskCompletion(logger *slog.Logger, taskId uint64, status bool, userId uint64) error {
	query := "UPDATE tasks SET is_completed = ? WHERE user_id = ? AND id = ?;"

	result, err := db.Exec(query, status, userId, taskId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec SetTaskCompletion", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected SetTaskCompletion", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected SetTaskCompletion", slog.String("err", "There were no rows affected"))
	}

	return nil
}
