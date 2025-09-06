package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type (
	ResponseStamp struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		CreatorID    uuid.UUID `json:"creatorId"`
		FileID       uuid.UUID `json:"fileId"`
		IsUnicode    bool      `json:"isUnicode"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
		HasThumbnail bool      `json:"hasThumbnail"`
	}

	StampData struct {
		ID        uuid.UUID `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		CreatorID uuid.UUID `db:"creator_id" json:"creator_id"`
		FileID    uuid.UUID `db:"file_id" json:"file_id"`
		IsUnicode bool      `db:"is_unicode" json:"is_unicode"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}
)

func (r *Repository) SaveStamp(ctx context.Context, stamps []*ResponseStamp) error {
	if len(stamps) == 0 {
		return nil
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var ids []uuid.UUID
	for _, s := range stamps {
		ids = append(ids, s.ID)
	}

	existingStamps, err := r.FindByID(ctx, tx, ids)
	if err != nil {
		return fmt.Errorf("failed to find exisingStamps: %w", err)
	}

	var inserts []*StampData
	var updates []*StampData

	for _, s := range stamps {
		if existingUpdatedAt, ok := existingStamps[s.ID]; ok {
			if !existingUpdatedAt.Equal(s.UpdatedAt) {
				updates = append(updates, &StampData{
					ID:        s.ID,
					Name:      s.Name,
					CreatorID: s.CreatorID,
					FileID:    s.FileID,
					IsUnicode: s.IsUnicode,
					CreatedAt: s.CreatedAt,
					UpdatedAt: s.UpdatedAt,
				})
			}
		} else {

			inserts = append(inserts, &StampData{
				ID:        s.ID,
				Name:      s.Name,
				CreatorID: s.CreatorID,
				FileID:    s.FileID,
				IsUnicode: s.IsUnicode,
				CreatedAt: s.CreatedAt,
				UpdatedAt: s.UpdatedAt,
			})
		}
	}
	if len(inserts) > 0 {
		if err := r.InsertStamps(ctx, tx, inserts); err != nil {
			return fmt.Errorf("failed to insert stamps: %w", err)
		}
	}
	if len(updates) > 0 {
		if err := r.UpdateStamps(ctx, tx, updates); err != nil {
			return fmt.Errorf("failed to update stamps: %w", err)
		}
	}

	return tx.Commit()

}

func (r *Repository) FindByID(ctx context.Context, tx *sqlx.Tx, ids []uuid.UUID) (map[uuid.UUID]time.Time, error) {

	query, args, err := sqlx.In("SELECT id, updated_at FROM stamps WHERE id IN (?)", ids)

	if err != nil {
		return nil, fmt.Errorf("failed to create IN query: %w", err)
	}
	query = tx.Rebind(query)

	updatedAts := make(map[uuid.UUID]time.Time)

	rows, err := tx.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to scan stamp row: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var updatedAt time.Time
		if err := rows.Scan(&id, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stamp row: %w", err)
		}
		updatedAts[id] = updatedAt
	}

	return updatedAts, nil
}

func (r *Repository) InsertStamps(ctx context.Context, tx *sqlx.Tx, stamps []*StampData) error {
	_, err := tx.NamedExecContext(ctx, `
        INSERT INTO stamps(id, name, creator_id, file_id, is_unicode, created_at, updated_at, count_monthly, count_total)
        VALUES (:id, :name, :creator_id , :file_id, :is_unicode, :created_at, :updated_at, 0, 0)
    `, stamps)

	if err != nil {
		log.Printf("Error inserting stamps: %v", err)

		return fmt.Errorf("failed to bulk insert stamps: %w", err)
	}

	return nil
}

func (r *Repository) UpdateStamps(ctx context.Context, tx *sqlx.Tx, stamps []*StampData) error {
	for _, s := range stamps {
		_, err := tx.NamedExecContext(ctx, `
            UPDATE stamps SET
                name = :name,
                creator_id = :creator_id,
                file_id = :file_id,
                is_unicode = :is_unicode,
                updated_at = :updated_at
            WHERE id = :id
        `, s)
		if err != nil {
			return fmt.Errorf("failed to update stamp %s: %w", s.ID, err)
		}
	}

	return nil
}
