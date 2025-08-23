package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) IsAdmin(context context.Context, userID uuid.UUID) (any, error) {
	panic("unimplemented")
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}
