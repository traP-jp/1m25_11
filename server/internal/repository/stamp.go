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
		ID           uuid.UUID `db:"id" json:"id"`
		Name         string    `db:"name" json:"name"`
		FileID       uuid.UUID `db:"file_id" json:"file_id"`
		CreatorID    uuid.UUID `db:"creator_id" json:"creator_id"`
		IsUnicode    bool      `db:"is_unicode" json:"is_unicode"`
		CreatedAt    time.Time `db:"created_at" json:"created_at"`
		UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
		CountMonthly int       `db:"count_monthly" json:"count_monthly"`
		CountTotal   int64     `db:"count_total" json:"count_total"`
	}

	StampSummary struct {
		ID     uuid.UUID `db:"id" json:"id"`
		Name   string    `db:"name" json:"name"`
		FileID uuid.UUID `db:"file_id" json:"file_id"`
	}
)

func (r *Repository) GetStampDetails(ctx context.Context) ([]*Stamp, error) {
	stamps := []*Stamp{}
	if err := r.db.SelectContext(ctx, &stamps, "SELECT * FROM stamps"); err != nil {
		return nil, fmt.Errorf("select stamps: %w", err)
	}

	return stamps, nil
}

func (r *Repository) GetStampSummaries(ctx context.Context) ([]*StampSummary, error) {
	stampSummaries := []*StampSummary{}
	if err := r.db.SelectContext(ctx, &stampSummaries, "SELECT id,name,file_id FROM stamps"); err != nil {
		return nil, fmt.Errorf("select stamps: %w", err)
	}

	return stampSummaries, nil
}

func (r *Repository) GetStampsByTagID(ctx context.Context, tagID uuid.UUID) ([]*Stamp, error) {
	stampsByTagID := []*Stamp{}
	query := `SELECT
            stamps.id, stamps.name, stamps.file_id, stamps.creator_id,
            stamps.is_unicode, stamps.created_at, stamps.updated_at,
            stamps.count_monthly, stamps.count_total
        FROM stamps
        INNER JOIN stamp_tags ON stamps.id = stamp_tags.stamp_id
        WHERE stamp_tags.tag_id = ?`
	if err := r.db.SelectContext(ctx, &stampsByTagID, query, tagID); err != nil {
		return nil, fmt.Errorf("select stamps by tagID: %w", err)
	}

	return stampsByTagID, nil
}

func (r *Repository) GetStampByStampID(ctx context.Context, stampID uuid.UUID) (*Stamp, error) {
	stampByStampID := &Stamp{}
	if err := r.db.GetContext(ctx, stampByStampID, "SELECT * FROM stamps WHERE id = ?", stampID); err != nil {
		return nil, fmt.Errorf("select stamps by stampID: %w", err)
	}

	return stampByStampID, nil
}
