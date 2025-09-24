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
	Hash         string
	CreationTime time.Time
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type Category struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	UserId uint64 `json:"userId"`
}
