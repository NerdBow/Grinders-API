package service

import (
	"log/slog"
	"time"

	"github.com/NerdBow/Grinders-API/internal/auth"
	"github.com/NerdBow/Grinders-API/internal/database"
	"github.com/NerdBow/Grinders-API/internal/util"
)

const MIN_PASSWORD_LENGHT = 8

type AuthService struct {
	userDb        database.UsersDB
	sessionDb     database.SessionsDB
	authSettings  auth.ArgonSettings
	tokenSettings auth.TokenSettings
}

func NewAuthService(userDb database.UsersDB, sessionDb database.SessionsDB, authSettings auth.ArgonSettings, tokenSettings auth.TokenSettings) AuthService {
	return AuthService{
		userDb:        userDb,
		sessionDb:     sessionDb,
		authSettings:  authSettings,
		tokenSettings: tokenSettings,
	}
}

func (s *AuthService) RegisterNewUser(username string, password string) error {
	if username == "" {
		slog.Info("")
		return nil
	}
	if len(password) < MIN_PASSWORD_LENGHT {
		slog.Info("")
		return nil
	}

	user := util.User{
		Username:     username,
		Hash:         s.authSettings.CreateNewHash(password),
		CreationTime: time.Now().UTC(),
	}
	err := s.userDb.AddUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(username string, password string) (util.Tokens, error) {
	user, err := s.userDb.GetUserByUsername(username)
	if err != nil {
		return util.Tokens{}, err
	}

	isValid := s.authSettings.CompareHash(password, user.Hash)
	if !isValid {
		slog.Info("")
		return util.Tokens{}, err
	}

	access, _, err := s.tokenSettings.CreateAccessToken(user.Id)
	if err != nil {
		return util.Tokens{}, err
	}

	refresh, err := s.tokenSettings.CreateRefreshToken(user.Id)
	if err != nil {
		return util.Tokens{}, err
	}

	session := util.Session{
		HashedId:       auth.HashRefreshTokenId(refresh.Id),
		ExpirationTime: refresh.Expire,
		CreationTime:   refresh.CreationTime,
		UserId:         refresh.UserId,
	}
	s.sessionDb.AddSession(session)

	tokens := util.Tokens{
		Access:  access,
		Refresh: refresh.Id,
	}

	return tokens, nil
}

func (s *AuthService) Refresh(refreshToken string, userId uint64) (util.Tokens, error) {
	if userId < 1 {
		return util.Tokens{}, nil // REPLACE
	}

	session, err := s.sessionDb.GetSession(auth.HashRefreshTokenId(refreshToken), userId)

	if err != nil {
		return util.Tokens{}, err // TODO change erorr
	}
	if session.UserId == 0 {
		return util.Tokens{}, err // TODO change erorr
	}
	if session.HashedId == "" {
		return util.Tokens{}, err // TODO change erorr
	}
	if time.Now().UTC().After(session.ExpirationTime) {
		return util.Tokens{}, util.ErrSessionExpired // TODO change error
	}

	access, _, err := s.tokenSettings.CreateAccessToken(userId)
	if err != nil {
		return util.Tokens{}, err
	}

	refresh, err := s.tokenSettings.CreateRefreshToken(userId)
	if err != nil {
		return util.Tokens{}, err
	}

	err = s.sessionDb.DeleteSession(session.HashedId)
	if err != nil {
		return util.Tokens{}, err
	}

	session = util.Session{
		HashedId:       auth.HashRefreshTokenId(refresh.Id),
		ExpirationTime: refresh.Expire,
		CreationTime:   refresh.CreationTime,
		UserId:         refresh.UserId,
	}
	err = s.sessionDb.AddSession(session)

	if err != nil {
		return util.Tokens{}, err
	}

	tokens := util.Tokens{
		Access: access,
		Refresh: refresh.Id,
	}

	return tokens, nil
}
