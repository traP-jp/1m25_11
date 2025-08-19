package repository

import (
	"context"
	"fmt"
)

type StampCountResult struct {
	StampID string `db:"stamp_id"`
	Count   int    `db:"total_count"`
}

func (r *Repository) GetStampCount(ctx context.Context, Since *openapi_types.Date, Until *openapi_types.Date) ([]StampCountResult, error) {
	results := []StampCountResult{}
	query := `
		SELECT
			stamp_id,
			SUM(reaction_count + message_count) AS total_count
		FROM
			stamp_daily_usages
		WHERE
			date BETWEEN ? AND ?
		GROUP BY
			stamp_id
		ORDER BY
			total_count DESC;
	`
	if err := r.db.SelectContext(ctx, &results, query, Since, Until); err != nil {
		return nil, fmt.Errorf("failed to get stamp count: %w", err)
	}

	return results, nil
}
