package domain

import (
	"time"
)

type Status string

const (
	InProgress Status = "in_progress"
	Close      Status = "close"
)

func (s Status) String() string {
	return string(s)
}

type Reception struct {
	ID       string
	DateTime time.Time
	PVZID    string
	Status   Status
}
