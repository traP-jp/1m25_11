package repository

import (
	"context"

	"github.com/google/uuid"
)

type (
	StampTag struct {
		id        uuid.UUID      `db:"id"`
		name      string         `db:"name"`
		createrID uuid.UUID      `db:"creater_id"`
		createdAt string         `db:"created_at"`
		updatedAt string         `db:"updated_at"`
		count     int            `db:"count"`
		stamps    []stampSummary `db:"stamps"`
	}

	CreateStampTagsParams struct {
		name      string
		createrID uuid.UUID
		createdAt string
		updatedAt string
	}
)

func (r *Repository) CreateStampTags(ctx context.Context, params CreateStampTagsParams) (uuid.UUID, error) {
	stampTagID := uuid.New()
	if _, err := r.db.ExecContext(ctx, "INSERT INTO stamp_tags (id, name, creater_id) VALUES (?, ?, ?)", stampTagID, params.name, params.createrID, params.createdAt, params.updatedAt); err != nil {
		return uuid.Nil, err
	}

	return stampTagID, nil
}
