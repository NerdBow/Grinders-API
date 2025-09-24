package sqlite

import (
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

func (db *SQLiteDB) AddCategory(name string, userId uint64) error {
	query := "INSERT INTO categories (name, user_id) VALUES (?, ?);"
	result, err := db.Exec(query, name, userId)
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

func (db *SQLiteDB) GetCategory(name string, userId uint64) (util.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE user_id=? AND name=?;"
	row := db.QueryRow(query, userId, name)

	category := util.Category{}
	err := row.Scan(&category.Id, &category.Name, &category.UserId)
	if err != nil {
		slog.Error("")
		return category, err
	}
	return category, nil
}

func (db *SQLiteDB) QueryCategory(prefix string, userId uint64) ([]util.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE user_id=? AND name LIKE ?;"
	rows, err := db.Query(query, userId, prefix + "%")
	if err != nil {
		slog.Error("")
		return nil, err
	}

	categories := make([]util.Category, 0, 10)
	for rows.Next() {
		category := util.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.UserId)
		if err != nil {
			slog.Error("")
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *SQLiteDB) GetUserCategories(userId uint64) ([]util.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE user_id=? ORDER BY name ASC;"
	rows, err := db.Query(query, userId)
	if err != nil {
		slog.Error("")
		return nil, err
	}

	categories := make([]util.Category, 0, 10)
	for rows.Next() {
		category := util.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.UserId)
		if err != nil {
			slog.Error("")
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *SQLiteDB) EditCategoryName(categoryId uint64, newName string, userId uint64) error {
	query := "UPDATE categories SET name = ? WHERE user_id = ? AND id = ?;"

	result, err := db.Exec(query, newName, userId, categoryId)
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

func (db *SQLiteDB) DeleteCategory(categoryId uint64, userId uint64) error {
	query := "DELETE FROM categories WHERE user_id = ? AND id = ?;"

	result, err := db.Exec(query, userId, categoryId)
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
