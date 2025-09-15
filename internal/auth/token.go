package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrKeyNotFound       = errors.New("Unable to get JWT_SIGNING_KEY from the env variables.")
)

type AccessToken struct {
	Sub uint64
	Aud []string
	Exp time.Time
	Iat time.Time
	Jti string
}

func CreateAccessToken(userId uint64) (string, error) {
	key := os.Getenv("JWT_SIGNING_KEY")

	tokenId := uuid.New().String()
	if tokenId == "" {
		return "", ErrKeyNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "", // I don't have a domain, so I will just leave this blank
		Subject:   strconv.FormatUint(userId, 10),
		Audience:  []string{"GrindersTUI"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * 5)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ID:        tokenId,
	})

	s, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return s, nil
}

func ParseAccessToken(tokenString string) (AccessToken, error) {
	keyFunc := func(t *jwt.Token) (any, error) {
		key := os.Getenv("JWT_SIGNING_KEY")
		if key == "" {
			return nil, ErrKeyNotFound
		}
		return []byte(key), nil
	}

	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return AccessToken{}, err
	}

	accessToken := AccessToken{}

	accessToken.Sub, err = strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return AccessToken{}, fmt.Errorf("Unable cast subject string into uint64: %w", err)
	}

	accessToken.Aud = claims.Audience
	accessToken.Exp = claims.ExpiresAt.Time
	accessToken.Iat = claims.IssuedAt.Time
	accessToken.Jti = claims.ID

	slog.Debug("Access Token", slog.Any("accessToken", accessToken))

	return accessToken, nil
}
