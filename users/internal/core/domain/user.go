package domain

import (
	"time"
)

type User struct {
	ID        int64
	Email     string
	CreateAt  time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
