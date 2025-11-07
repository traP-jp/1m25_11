package repository

import (
	"context"
	"github.com/google/uuid"
	"log"
	"strings"
)

type (
	StampStatus struct {
		ID         uuid.UUID `db:"id"`
		TotalCount int       `db:"total_count"`
		Count      int       `db:"count"`
	}
)

func (r *Repository) updateUsage() error {

	return nil
}

func (r *Repository) UpdateCount(ctx context.Context, stampTotalCount map[uuid.UUID]int, rawCCount map[uuid.UUID]int) error {
	log.Printf("Updating total counts for %d stamps", len(stampTotalCount))
	tx, err := r.db.BeginTxx(ctx, nil)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err != nil {
		log.Printf("Error starting transaction: ", err)
		return err
	}
	var caseBuilder strings.Builder
	args := make([]interface{}, 0, len(stampTotalCount)*2)
	ids := make([]interface{}, 0, len(stampTotalCount))
	for id, totalCount := range stampTotalCount {
		caseBuilder.WriteString("WHEN ? THEN ? ")
		args = append(args, id, totalCount)
		ids = append(ids, id)
	}
	query := `UPDATE stamps SET count_total = CASE id ` + caseBuilder.String() + `ELSE count_total END`
	log.Print(query, args, ids)

	query = r.db.Rebind(query)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing update: ", err)
		return err
	}

	log.Printf("Successfully updated total counts for %d stamps", len(stampTotalCount))

	return nil
}
