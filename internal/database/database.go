package database

import "github.com/NerdBow/Grinders-API/internal/util"

type SessionsDB interface {
	// AddSession will inserts the given session information into the database.
	AddSession(session util.Session) error
	// GetSession will query a session with the specified hashedId and userId.
	// If there is no session then an empty Session will be returned.
	GetSession(hashedId string, userId uint64) (util.Session, error)
	// DeleteSesssion will delete a session specified by the given hashedId.
	// If unable to delete the session an error will be returned.
	DeleteSession(hashedId string) error
}

type UsersDB interface {
	// AddUser will insert the given user information into the database.
	AddUser(user util.User) error
	// GetUser will query a user with the specified userId.
	// If there is no user with the userId, then an error will be returned.
	GetUser(userId uint64) (util.User, error)
	// EditUsername will change the username of the user specified by userId to the given newName.
	// If there is no user with the userId, then no error will be returned.
	// Errors are only returned for database errors.
	EditUsername(userId uint64, newName string) error
}

type (
	TasksDB        interface{}
	WorkLogsDB     interface{}
	PausesDB       interface{}
	BreaksDB       interface{}
	CategoriesDB   interface{}
	GoalsDB        interface{}
	GroupsDB       interface{}
	GroupMembersDB interface{}
)
