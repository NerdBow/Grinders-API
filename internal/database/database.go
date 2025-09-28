package database

import (
	"log/slog"

	"github.com/NerdBow/Grinders-API/internal/util"
)

type SessionsDB interface {
	// AddSession will inserts the given session information into the database.
	AddSession(logger *slog.Logger, session util.Session) error
	// GetSession will query a session with the specified hashedId and userId.
	// If there is no session then an empty Session will be returned.
	GetSession(logger *slog.Logger, hashedId string, userId uint64) (util.Session, error)
	// DeleteSesssion will delete a session specified by the given hashedId.
	// If unable to delete the session an error will be returned.
	DeleteSession(logger *slog.Logger, hashedId string) error
}

type UsersDB interface {
	// AddUser will insert the given user information into the database.
	AddUser(logger *slog.Logger, user util.User) error
	// GetUser will query a user with the specified userId.
	// If there is no user with the userId, then an error will be returned.
	GetUser(logger *slog.Logger, userId uint64) (util.User, error)
	// GetUserByUsername will query a user with the specified username.
	// If there is no user with the userId, then an error will be returned.
	GetUserByUsername(logger *slog.Logger, username string) (util.User, error)
	// EditUsername will change the username of the user specified by userId to the given newName.
	// If there is no user with the userId, then no error will be returned.
	// Errors are only returned for database errors.
	EditUsername(logger *slog.Logger, userId uint64, newName string) error
}

type CategoriesDB   interface{
	// AddCategory will create a new category with the specified name for the userId.
	AddCategory(logger *slog.Logger, name string, userId uint64) error
	// GetCategory will retrive the specific category specified by name.
	GetCategory(logger *slog.Logger, name string, userId uint64) (util.Category, error)
	// QueryCategory will retrive ALL categories prefixed with the specified prefix.
	QueryCategory(logger *slog.Logger, prefix string, userId uint64) ([]util.Category, error)
	// GetAllUserCategories will retrive all categories linked to the userId.
	// The slice of Category structs will be sorted in alphabetical order by category name.
	GetUserCategories(logger *slog.Logger, userId uint64) ([]util.Category, error)
	// EditCategoryName will change the name of the category for categoryId to newName.
	EditCategoryName(logger *slog.Logger, categoryId uint64, newName string, userId uint64) error
	// DeleteCategory will delete the category with the specified categoryId.
	DeleteCategory(logger *slog.Logger, categoryId uint64, userId uint64) error
}

type (
	TasksDB        interface{}
	WorkLogsDB     interface{}
	PausesDB       interface{}
	BreaksDB       interface{}
	GoalsDB        interface{}
	GroupsDB       interface{}
	GroupMembersDB interface{}
)
