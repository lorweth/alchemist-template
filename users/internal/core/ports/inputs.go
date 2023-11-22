package ports

import (
	"time"
)

type GetUserInput struct {
	ID        int64
	Email     string
	CreatedAt Period
}

type Period struct {
	From time.Time
	To   time.Time
}
