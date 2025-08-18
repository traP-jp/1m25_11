package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	// stamps table
	Stamp struct {
		ID    uuid.UUID `db:"id"`
		Name  string    `db:"name"`
		FileID uuid.UUID    `db:"file_id"`
		CreatorID uuid.UUID `db:"creator_id"`
		IsUnicode bool `db:"is_unicode"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		CountMonthly int `db:"count_monthly"`
		CountTotal int64 `db:"count_total"`
	}

)

func (r *Repository) GetStamps(ctx context.Context) ([]*Stamp, error) {
	stamps := []*Stamp{}
	if err := r.db.SelectContext(ctx, &stamps, "SELECT * FROM stamps"); err != nil {
		return nil, fmt.Errorf("select stamps: %w", err)
	}
	
	return stamps, nil
}

func (r *Repository) GetStampsByTagID(ctx context.Context, tagID uuid.UUID) ([]*Stamp, error) {
	stampsByTagID := []*Stamp{}
	if err := r.db.SelectContext(ctx, &stampsByTagID, "SELECT * FROM stamps JOIN stamp_tags ON stamps.id = stamp_tags.stamp_id WHERE stamp_tags.tag_id = ?", tagID); err != nil {
		return nil, fmt.Errorf("select stamps by tagID: %w", err)
	}

	return stampsByTagID, nil
}

func (r *Repository) GetStampsByStampID(ctx context.Context, stampID uuid.UUID) (*Stamp, error) {
	stampsByStampID := &Stamp{}
	if err := r.db.GetContext(ctx, stampsByStampID, "SELECT * FROM stamps WHERE id = ?", stampID); err != nil {
		return nil, fmt.Errorf("select stamps by stampID: %w", err)
	}

	return stampsByStampID, nil
}