package util

import "time"

const (
	SORT_CREATION uint8 = iota + 1 // Reserve 0 for no op
	SORT_COMPLETION
	SORT_DEADLINE
	ORDER_ASCEDNING
	ORDER_DESCEDNING
)

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

type Task struct {
	Id             uint64
	Name           string
	CreationTime   time.Time
	CompletionTime time.Time
	DeadlineTime   time.Time
	IsComplete     bool
	CategoryId     uint64
	UserId         uint64
}

type TaskQuerySettings struct {
	Name      string
	Category  string
	SortType  uint8
	SortOrder uint8
	Page      uint16
	UserId    uint64
}
