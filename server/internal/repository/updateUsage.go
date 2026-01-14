package repository

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
)

func (r *Repository) UpdateCount(ctx context.Context, stampTotalCount map[uuid.UUID]int, rawCount map[uuid.UUID]int) error {
	if len(stampTotalCount) == 0 {
		log.Printf("UpdateCount: No stamps to update.")
		return nil
	}
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
		log.Printf("Error starting transaction:%v ", err)
		return err
	}
	var caseBuilderTotal strings.Builder
	var caseBuilderRaw strings.Builder
	argsTotal := make([]interface{}, 0, len(stampTotalCount)*2)
	argsRaw := make([]interface{}, 0, len(stampTotalCount)*2)
	idsTotal := make([]interface{}, 0, len(stampTotalCount))
	idsRaw := make([]interface{}, 0, len(stampTotalCount))
	for id, totalCount := range stampTotalCount {
		rawC, ok := rawCount[id]
		if !ok {
			rawC = 0
		}
		caseBuilderTotal.WriteString("WHEN ? THEN ? ")
		caseBuilderRaw.WriteString("WHEN ? THEN ? ")
		argsTotal = append(argsTotal, id, totalCount)
		argsRaw = append(argsRaw, id, rawC)
		idsTotal = append(idsTotal, id)
		idsRaw = append(idsRaw, id)
	}
	queryTotal := `UPDATE stamps SET count_total = CASE id ` + caseBuilderTotal.String() + `ELSE count_total END`
	queryRaw := `UPDATE stamps SET count = CASE id ` + caseBuilderRaw.String() + `ELSE count END`
	

	queryTotal = r.db.Rebind(queryTotal)
	queryRaw = r.db.Rebind(queryRaw)

	_, err = tx.ExecContext(ctx, queryTotal, argsTotal...)
	if err != nil {
		log.Printf("Error executing update:%v ", err)

		return err
	}
	_, err = tx.ExecContext(ctx, queryRaw, argsRaw...)
	if err != nil {
		log.Printf("Error executing update:%v ", err)

		return err
	}
	log.Printf("Successfully updated both counts")

	return nil
}
