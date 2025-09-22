package database

import "github.com/NerdBow/Grinders-API/internal/util"

type SessionsDB interface{
	// AddSession will inserts the given session information into the database.
	AddSession(session util.Session) error
	// GetSession will query a session with the specified hashedId and userId.
	// If there is no session then an empty Session will be returned.
	GetSession(hashedId string, userId uint64) (util.Session, error)
	// DeleteSesssion will delete a session specified by the given hashedId.
	// If unable to delete the session an error will be returned.
	DeleteSession(hashedId string) error
}
type UsersDB interface{}
type TasksDB interface{}
type WorkLogsDB interface{}
type PausesDB interface{}
type BreaksDB interface{}
type CategoriesDB interface{}
type GoalsDB interface{}
type GroupsDB interface{}
type GroupMembersDB interface{}
