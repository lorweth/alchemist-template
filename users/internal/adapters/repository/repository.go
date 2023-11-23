package repository

import (
	"github.com/virsavik/alchemist-template/pkg/postgres"
)

type Repository struct {
	db postgres.ContextExecutor
}

func New(db postgres.ContextExecutor) *Repository {
	return &Repository{
		db: db,
	}
}
