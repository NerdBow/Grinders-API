package util

import ("errors")

var (
	ErrInvalidUserId = errors.New("Invalid user id")
	ErrSessionExpired = errors.New("Session has expired")
)

