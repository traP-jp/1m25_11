package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StampRankingResult struct {
	StampID       uuid.UUID `db:"stamp_id"`
	Name          string    `db:"name"`
	FileID        uuid.UUID `db:"file_id"`
	ReactionCount int       `db:"reaction_count"`
	MessageCount  int       `db:"message_count"`
}

func (r *Repository) GetRanking(ctx context.Context, since, until time.Time) ([]StampRankingResult, error) {
	var results []StampRankingResult
	sinceDate := since.Truncate(24 * time.Hour)
	untilDate := until.Truncate(24 * time.Hour)
	query := `
        SELECT
            s.id AS stamp_id,
            s.name,
            s.file_id,
            SUM(sdu.reaction_count) AS reaction_count,
            SUM(sdu.message_count) AS message_count
        FROM
            stamps s
        INNER JOIN
            stamp_daily_usages sdu ON s.id = sdu.stamp_id
        WHERE
            sdu.date BETWEEN ? AND ?
        GROUP BY
            s.id, s.name, s.file_id
        HAVING
            SUM(sdu.reaction_count) > 0 OR SUM(sdu.message_count) > 0
        ORDER BY
            reaction_count DESC, message_count DESC, s.name ASC;
    `
	err := r.db.SelectContext(ctx, &results, query, sinceDate, untilDate)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return []StampRankingResult{}, nil
	}

	return results, nil
}

