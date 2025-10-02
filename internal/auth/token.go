package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrUnableToGenerateUUID = errors.New("Unable to generate UUID for token.")

type JWTToken struct {
	Sub uint64
	Aud []string
	Exp time.Time
	Iat time.Time
	Jti string
}

type RefreshToken struct {
	Id           string
	Expire       time.Time
	CreationTime time.Time
	UserId       uint64
}

type TokenSettings struct {
	AccessDuration  time.Duration // In minutes for env var ACCESS_TOKEN_DURATION
	SessionDuration time.Duration // In hours for env var SESSION_DURATION
}

// NewTokenSettings creates a new TokenSettings with its field specified by the related env vars.
func NewTokenSettings() TokenSettings {
	settings := TokenSettings{}
	s := os.Getenv("ACCESS_TOKEN_DURATION")
	n, _ := strconv.ParseUint(s, 10, 64)
	settings.AccessDuration = time.Duration(time.Minute * time.Duration(n))
	s = os.Getenv("SESSION_DURATION")
	n, _ = strconv.ParseUint(s, 10, 64)
	settings.SessionDuration = time.Duration(time.Hour * time.Duration(n))
	return settings
}

// CreateAccessToken creates and returns a JWT access token and its uuid for the given user id.
func (s *TokenSettings) CreateAccessToken(userId uint64) (string, string, error) {
	key := os.Getenv("JWT_SIGNING_KEY")

	tokenId := uuid.New().String()
	if tokenId == "" {
		slog.LogAttrs(context.Background(), slog.LevelError, "uuid.New CreateAccessToken", slog.String("err", ErrUnableToGenerateUUID.Error()))
		return "", "", ErrUnableToGenerateUUID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "", // I don't have a domain, so I will just leave this blank
		Subject:   strconv.FormatUint(userId, 10),
		Audience:  []string{"GrindersTUI"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(s.AccessDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ID:        tokenId,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "SignedString CreateAccessToken", slog.String("err", err.Error()))
		return "", "", err
	}
	return tokenString, tokenId, err
}

// CreateRefreshToken creates and returns a RefreshToken struct for the given user id and Access Token Id.
func (s *TokenSettings) CreateRefreshToken(userId uint64) (RefreshToken, error) {
	tokenId := uuid.New().String()
	if tokenId == "" {
		slog.LogAttrs(context.Background(), slog.LevelError, "uuid.New CreateRefreshToken", slog.String("err", ErrUnableToGenerateUUID.Error()))
		return RefreshToken{}, ErrUnableToGenerateUUID
	}

	return RefreshToken{
		Id:           tokenId,
		Expire:       time.Now().UTC().Add(s.SessionDuration),
		CreationTime: time.Now().UTC(),
		UserId:       userId,
	}, nil
}

// ParseToken parses a given JWT token string and returns a JWTToken of the parsed token claims.
func ParseToken(tokenString string) (JWTToken, error) {
	keyFunc := func(t *jwt.Token) (any, error) {
		key := os.Getenv("JWT_SIGNING_KEY")
		return []byte(key), nil
	}

	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelDebug, "ParseWithClaims ParseToken", slog.String("err", err.Error()))
		return JWTToken{}, err
	}

	jwtToken := JWTToken{}

	jwtToken.Sub, err = strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelDebug, "sub parsing ParseToken", slog.String("err", err.Error()))
		return JWTToken{}, fmt.Errorf("Unable cast subject string into uint64: %w", err)
	}

	jwtToken.Aud = claims.Audience
	jwtToken.Exp = claims.ExpiresAt.Time
	jwtToken.Iat = claims.IssuedAt.Time
	jwtToken.Jti = claims.ID

	slog.LogAttrs(context.Background(), slog.LevelDebug, "parsed token ParseToken", slog.Any("jwtStruct", jwtToken))

	return jwtToken, nil
}
