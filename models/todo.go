package models

import (
	"time"
)

type Task struct {
	ID        uint
	TaskName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
