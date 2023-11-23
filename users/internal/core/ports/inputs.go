package ports

import (
	"time"

	"github.com/virsavik/alchemist-template/users/internal/pkg/pagination"
)

type GetUserInput struct {
	ID         int64
	Email      string
	CreatedAt  Period
	Pagination pagination.Input
}

type Period struct {
	From time.Time
	To   time.Time
}
