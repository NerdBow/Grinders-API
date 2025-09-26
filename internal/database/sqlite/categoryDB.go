package sqlite

import (
	"context"
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddCategory(logger *slog.Logger, name string, userId uint64) error {
	query := "INSERT INTO categories (name, user_id) VALUES (?, ?);"
	result, err := db.Exec(query, name, userId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec AddCategory", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddCategory", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected AddCategory", slog.String("err", "There were no rows affected"))
	}

	return nil
}

func (db *SQLiteDB) GetCategory(logger *slog.Logger, name string, userId uint64) (util.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE user_id=? AND name=?;"
	row := db.QueryRow(query, userId, name)

	category := util.Category{}
	err := row.Scan(&category.Id, &category.Name, &category.UserId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Scan GetCategory", slog.String("err", err.Error()))
		return category, err
	}
	return category, nil
}

func (db *SQLiteDB) QueryCategory(logger *slog.Logger, prefix string, userId uint64) ([]util.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE user_id=? AND name LIKE ?;"
	rows, err := db.Query(query, userId, prefix+"%")
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Query QueryCategory", slog.String("err", err.Error()))
		return nil, err
	}

	categories := make([]util.Category, 0, 10)
	for rows.Next() {
		category := util.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.UserId)
		if err != nil {
			logger.LogAttrs(context.Background(), slog.LevelError, "Scan QueryCategory", slog.String("err", err.Error()))
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *SQLiteDB) GetUserCategories(logger *slog.Logger, userId uint64) ([]util.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE user_id=? ORDER BY name ASC;"
	rows, err := db.Query(query, userId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Query GetUserCategories", slog.String("err", err.Error()))
		return nil, util.ErrDatabase
	}

	categories := make([]util.Category, 0, 10)
	for rows.Next() {
		category := util.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.UserId)
		if err != nil {
			logger.LogAttrs(context.Background(), slog.LevelError, "Scan GetUserCategories", slog.String("err", err.Error()))
			return nil, util.ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *SQLiteDB) EditCategoryName(logger *slog.Logger, categoryId uint64, newName string, userId uint64) error {
	query := "UPDATE categories SET name = ? WHERE user_id = ? AND id = ?;"

	result, err := db.Exec(query, newName, userId, categoryId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec EditCategoryName", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected EditCategoryName", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected EditCategoryName", slog.String("err", "There were no rows affected"))
	}

	return nil
}

func (db *SQLiteDB) DeleteCategory(logger *slog.Logger, categoryId uint64, userId uint64) error {
	query := "DELETE FROM categories WHERE user_id = ? AND id = ?;"

	result, err := db.Exec(query, userId, categoryId)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Exec DeleteCategory", slog.String("err", err.Error()))
		return util.ErrDatabase
	}

	n, err := result.RowsAffected()
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected DeleteCategory", slog.String("err", err.Error()))
	}

	if n != 1 {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "RowsAffected DeleteCategory", slog.String("err", "There were no rows affected"))
	}

	return nil
}
