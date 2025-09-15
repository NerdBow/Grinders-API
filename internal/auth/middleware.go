package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

const MIN_TOKEN_LENGTH = 30 // Mainly here so I don't get a string slicing error

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < MIN_TOKEN_LENGTH {
			m(w, http.StatusUnauthorized, "Unauthorized request", "No Bearer token")
			return
		} else if authHeader[:6] != "Bearer" {
			m(w, http.StatusUnauthorized, "Unauthorized request", "No Bearer token")
			return
		}
		tokenString := authHeader[7:]
		accessToken, err := ParseToken(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				m(w, http.StatusUnauthorized, "Unauthorized request", "Token is expired")
				return
			default:
				m(w, http.StatusBadRequest, "Unauthorized request", "Token is malformed")
				return
			}
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", accessToken.Sub)
		h(w, r.WithContext(ctx))
	}
}

// This function is temperary probably going to move it into a different package
func m(w http.ResponseWriter, statusCode int, error string, message string) {
	response := struct {
		Error string `json:"error"`
		Message string `json:"message"`
	}{error, message}
	responseBytes, _ := json.Marshal(response)
	w.WriteHeader(statusCode)
	_, err := w.Write(responseBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
