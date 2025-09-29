package util

import (
	"errors"
)

var (
	ErrHashMismatch      = errors.New("Hash of given password and account hash is mismatched")
	ErrBadPassword       = errors.New("Password should be at least 8 characters")
	ErrEmptyString       = errors.New("An empty string is invalid")
	ErrInvalidUserId     = errors.New("Invalid user id")
	ErrInvalidCategoryId = errors.New("Invalid category id")
	ErrSessionExpired    = errors.New("Session has expired")
	ErrDatabase          = errors.New("Database Error")
)
