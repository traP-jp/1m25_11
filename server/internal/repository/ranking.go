package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StampCountResult struct {
	StampID       uuid.UUID `db:"stamp_id"`
	ReactionCount int       `db:"reaction_count"`
	MessageCount  int       `db"message_count"`
}

func (r *Repository) GetStampCount(ctx context.Context, since time.Time, until time.Time) ([]StampCountResult, error) {
	results := []StampCountResult{}
	sinceDate := since.Truncate(24 * time.Hour)
	untilDate := until.Truncate(24 * time.Hour)
	query := `
		SELECT
			stamp_id,
			reaction_count,
			message_count
		FROM
			stamp_daily_usages
		WHERE
			date BETWEEN ? AND ?
		GROUP BY
			stamp_id;
	`
	if err := r.db.SelectContext(ctx, &results, query, sinceDate, untilDate); err != nil {
		return nil, fmt.Errorf("failed to get stamp count: %w", err)
	}

	return results, nil
}
