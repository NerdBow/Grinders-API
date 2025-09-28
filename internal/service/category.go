package service

import (
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/database"
	"github.com/NerdBow/Grinders-API/internal/util"
)

type CategoryService struct {
	categoryDb database.CategoriesDB
}

func NewCategoryService(categoryDb database.CategoriesDB) CategoryService {
	return CategoryService{
		categoryDb: categoryDb,
	}
}

func (s *CategoryService) CreateCategory(logger *slog.Logger, userId uint64, name string) error {
	if userId < 1 {
		return nil
	}
	if name == "" {
		return nil
	}

	err := s.categoryDb.AddCategory(logger, name, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) GetCategory(logger *slog.Logger, userId uint64, name string) (util.Category, error) {
	if userId < 1 {
		return util.Category{}, nil
	}
	if name == "" {
		return util.Category{}, nil
	}

	category, err := s.categoryDb.GetCategory(logger, name, userId)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *CategoryService) QueryCategory(logger *slog.Logger, userId uint64, prefix string) ([]util.Category, error) {
	if userId < 1 {
		return nil, nil
	}
	if prefix == "" {
		return nil, nil
	}

	categories, err := s.categoryDb.QueryCategory(logger, prefix, userId)
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *CategoryService) ChangeName(logger *slog.Logger, userId uint64, categoryId uint64, newName string) error {
	if userId < 1 {
		return nil
	}
	if categoryId < 1 {
		return nil
	}

	err := s.categoryDb.EditCategoryName(logger, categoryId, newName, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteCategory(logger *slog.Logger, userId uint64, categoryId uint64) error {
	if userId < 1 {
		return nil
	}
	if categoryId < 1 {
		return nil
	}

	err := s.categoryDb.DeleteCategory(logger, categoryId, userId)
	if err != nil {
		return err
	}

	return nil
}
