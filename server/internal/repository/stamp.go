package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"


)

type (
	// users table
	Stamp struct {
		ID    uuid.UUID `db:"id"`
		Name  string    `db:"name"`
		FileID string    `db:"file_id"`
		CreatorID uuid.UUID `db:"creator_id"`
		IsUnicode Boolean `db:"is_unicode"`
		createdAt  string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
		CountMonthly string `db:count_monthly`
		CountTotal string `db:cou`
	}

	CreateUserParams struct {
		Name  string
		Email string
	}
)

func (r *Repository) getStamps(ctx context.Context) ([]*User, error) {
	users := []*User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (r *Repository) getCertainStamps(ctx context.Context) ([]*User, error) {
	users := []*User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (r *Repository) getCertainStampDetails(ctx context.Context) ([]*User, error) {
	users := []*User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
} 