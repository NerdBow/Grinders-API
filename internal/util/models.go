package util

import "time"

type Session struct {
	HashedId       string
	ExpirationTime time.Time
	CreationTime   time.Time
	UserId         uint64
}

type User struct {
	Id           uint64
	Username     string
	Salt         string
	Hash         string
	CreationTime time.Time
}
