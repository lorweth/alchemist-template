package repository

import (
	"github.com/virsavik/alchemist-template/pkg/postgres"
)

type Repository struct {
	db postgres.ContextExecutor
}

func New(db postgres.ContextExecutor) *Repository {
	// Init id generator
	initIDGenerator()

	return &Repository{
		db: db,
	}
}
