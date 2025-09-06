package repository

import (
	"context"

	"github.com/google/uuid"
)

type StampRankingResult struct {
	StampID      uuid.UUID `db:"id"`
	TotalCount   int       `db:"count_total"`
	MonthlyCount int       `db:"count_monthly"`
}

func (r *Repository) GetRanking(ctx context.Context) ([]StampRankingResult, error) {
	var results []StampRankingResult
	query := `
        SELECT
            id,count_total,count_monthly FROM stamps`
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return []StampRankingResult{}, nil
	}

	return results, nil
}
