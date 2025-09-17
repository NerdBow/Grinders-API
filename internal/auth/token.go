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

const (
	ACCESS_TOKEN_DURATION  = time.Minute * 5
	REFRESH_TOKEN_DURATION = time.Hour * 24 * 7
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
	Id            string
	AccessTokenId string
	Expire        time.Time
	CreationTime  time.Time
	UserId        uint64
}

// CreateAccessToken creates and returns a JWT access token and its uuid for the given user id.
func CreateAccessToken(userId uint64) (string, string, error) {
	key := os.Getenv("JWT_SIGNING_KEY")

	tokenId := uuid.New().String()
	if tokenId == "" {
		return "", "", ErrUnableToGenerateUUID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "", // I don't have a domain, so I will just leave this blank
		Subject:   strconv.FormatUint(userId, 10),
		Audience:  []string{"GrindersTUI"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ACCESS_TOKEN_DURATION)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ID:        tokenId,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", "", err
	}
	return tokenString, tokenId, err
}

// CreateRefreshToken creates and returns a RefreshToken struct for the given user id and Access Token Id.
func CreateRefreshToken(userId uint64, accessTokenId string) (RefreshToken, error) {
	tokenId := uuid.New().String()
	if tokenId == "" {
		return RefreshToken{}, ErrUnableToGenerateUUID
	}

	return RefreshToken{
		Id: tokenId,
		AccessTokenId: accessTokenId,
		Expire: time.Now().UTC().Add(REFRESH_TOKEN_DURATION),
		CreationTime: time.Now().UTC(),
		UserId: userId,
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
		return JWTToken{}, err
	}

	jwtToken := JWTToken{}

	jwtToken.Sub, err = strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return JWTToken{}, fmt.Errorf("Unable cast subject string into uint64: %w", err)
	}

	jwtToken.Aud = claims.Audience
	jwtToken.Exp = claims.ExpiresAt.Time
	jwtToken.Iat = claims.IssuedAt.Time
	jwtToken.Jti = claims.ID

	slog.Debug("Token", slog.Any("jwtStruct", jwtToken))

	return jwtToken, nil
}
