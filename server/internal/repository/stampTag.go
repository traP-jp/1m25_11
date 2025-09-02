package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	CreateStampTagParams struct {
		StampID   uuid.UUID `db:"stamp_id" json:"stamp_id"`
		TagID     uuid.UUID `db:"tag_id" json:"tag_id"`
		CreatorID uuid.UUID `db:"creator_id" json:"creator_id"`
	}
	StampTagSummary struct {
		StampID uuid.UUID
		TagID   uuid.UUID
	}
)

func (r *Repository) CreateStampTags(ctx context.Context, params CreateStampTagParams) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO stamp_tags (stamp_id, tag_id, creator_id) VALUES (?, ?, ?)", params.StampID, params.TagID, params.CreatorID); err != nil {
		return fmt.Errorf("failed to insert stampTags:%w", err)
	}

	return nil
}

func (r *Repository) DeleteStampTags(ctx context.Context, stampID uuid.UUID, tagID uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM stamp_tags WHERE stamp_id = ? AND tag_id = ?", stampID, tagID); err != nil {
		return fmt.Errorf("failed to delete stampTag:%w", err)
	}

	return nil
}

func (r *Repository) GetSearchStampTags(ctx context.Context, keyword string) ([]uuid.UUID, error) {
	stampIDs := []uuid.UUID{}
	if err := r.db.SelectContext(ctx, &stampIDs, "SELECT DISTINCT stamp_id FROM stamp_tags JOIN tags ON stamp_tags.tag_id = tags.id WHERE tags.name LIKE ?", "%"+keyword+"%"); err != nil {
		return nil, fmt.Errorf("failed to get stampID by tagKeyword : %w", err)
	}

	return stampIDs, nil
}

func (r *Repository) getStampTagsByCreatorID(ctx context.Context, creatorID uuid.UUID) ([]StampTagSummary, error) {
	stampTags := []StampTagSummary{}
	if err := r.db.SelectContext(ctx, &stampTags, "SELECT * FROM stamp_tags WHERE creator_id = ?", creatorID); err != nil {
		return nil, fmt.Errorf("select stamp_tags by userID: %w", err)
	}

	return stampTags, nil
}
