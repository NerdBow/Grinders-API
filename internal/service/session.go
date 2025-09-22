package service

import (
	"time"

	"github.com/NerdBow/Grinders-API/internal/auth"
	"github.com/NerdBow/Grinders-API/internal/database"
	"github.com/NerdBow/Grinders-API/internal/util"
)

type SessionService struct {
	sessionDb     database.SessionsDB
	tokenSettings auth.TokenSettings
}

func NewSessionService(sessionDb database.SessionsDB, tokenSettings auth.TokenSettings) SessionService {
	return SessionService{
		sessionDb:     sessionDb,
		tokenSettings: tokenSettings,
	}
}

// CreateSession will create a new session for the specified userId and return refresh string id.
func (s *SessionService) CreateSession(userId uint64) (string, error) {
	if userId <= 0 {
		return "", util.ErrInvalidUserId
	}

	token, err := s.tokenSettings.CreateRefreshToken(userId)
	if err != nil {
		return "", err
	}

	session := util.Session{
		HashedId:       auth.HashRefreshTokenId(token.Id),
		ExpirationTime: token.Expire,
		CreationTime:   token.CreationTime,
		UserId:         token.UserId,
	}
	err = s.sessionDb.AddSession(session)
	if err != nil {
		return "", err
	}

	return token.Id, nil
}

// ValidiateSession checks if the given sessionId is in the session database and also check if the session is expired or not.
func (s *SessionService) ValidiateSession(sessionId string, userId uint64) (bool, error) {
	if userId <= 0 {
		return false, util.ErrInvalidUserId
	}

	hashedId := auth.HashRefreshTokenId(sessionId)
	session, err := s.sessionDb.GetSession(hashedId, userId)
	if err != nil {
		return false, err // TODO change erorr
	}
	if session.UserId == 0 {
		return false, nil
	}
	if session.HashedId == "" {
		return false, nil
	}
	if time.Now().UTC().After(session.ExpirationTime) {
		return false, util.ErrSessionExpired // TODO change error
	}

	return true, nil
}

// DeleteSession will delete the session with the specified sessionId.
func (s *SessionService) DeleteSession(sessionId string) error {
	hashedId := auth.HashRefreshTokenId(sessionId)
	err := s.sessionDb.DeleteSession(hashedId)
	if err != nil {
		return err
	}
	return nil
}
