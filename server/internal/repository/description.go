package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	CreateDescriptionParams struct {
		StampID     uuid.UUID `db:"stamp_id"`
		Description string    `db:"description"`
		CreatorID   uuid.UUID `db:"creator_id"`
		CreatedAt   time.Time `db:"created_at"`
	}
	Description struct {
		StampID     uuid.UUID `db:"stamp_id"`
		Description string    `db:"description"`
		CreatorID   uuid.UUID `db:"creator_id"`
		CreatedAt   time.Time `db:"created_at"`
	}
)

func (r *Repository) CreateDescriptions(ctx context.Context, params CreateDescriptionParams) error{
	if _, err := r.db.ExecContext(ctx, "INSERT INTO stamp_description_revisions(stamp_id,description,creator_id,created_at) VALUES(?,?,?,?)", params.StampID, params.Description, params.CreatorID, params.CreatedAt); err != nil {
		return fmt.Errorf("failed to insert description: %w", err)
	}

	return  nil
}

func (r *Repository) GetDescriptionsByStampID(ctx context.Context, stampID uuid.UUID) ([]*Description, error) {
	descriptions := []*Description{}
	if err := r.db.SelectContext(ctx, &descriptions, "SELECT * FROM stamp_description_revisions WHERE stamp_id = ?", stampID); err != nil {
		return nil, fmt.Errorf("failed to get descriptions by stampID: %w", err)
	}

	return descriptions, nil
}

func (r *Repository) DeleteDescriptions(ctx context.Context, stampID uuid.UUID, creatorID uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM stamp_description_revisions WHERE stamp_id = ? AND creator_id = ?", stampID, creatorID); err != nil {
		return fmt.Errorf("failed to delete description: %w", err)
	}

	return nil
}

func (r *Repository) UpdateDescriptions(ctx context.Context, stampID uuid.UUID, creatorID uuid.UUID, description string) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE stamp_description_revisions SET description = ? WHERE stamp_id = ? AND creator_id =?", description, stampID, creatorID); err != nil {
		return fmt.Errorf("failed to update description: %w", err)
	}

	return nil
}
