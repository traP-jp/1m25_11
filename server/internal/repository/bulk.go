package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type TagInfo struct {
	Name string
}

type CreatedTagInfo struct {
	ID   uuid.UUID
	Name string
}

type tagInsertData struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatorID uuid.UUID `db:"creator_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type StampMetaAddition struct {
	StampID     uuid.UUID   `json:"stamp_id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	Description string      `json:"description"`
}

type stampTagLinkData struct {
	StampID   uuid.UUID `db:"stamp_id" json:"stamp_id"`
	TagID     uuid.UUID `db:"tag_id" json:"tag_id"`
	CreatorID uuid.UUID `db:"creator_id" json:"creator_id"`
}

type stampDescriptionData struct {
	StampID     uuid.UUID `db:"stamp_id" json:"stamp_id"`
	Description string    `db:"description" json:"description"`
	CreatorID   uuid.UUID `db:"creator_id" json:"creator_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (r *Repository) BulkCreateTags(ctx context.Context, tags []TagInfo) ([]CreatedTagInfo, error) {
	if len(tags) == 0 {
		return []CreatedTagInfo{}, nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	creatorID, err := uuid.Parse("3b261ff3-f940-4e2c-a626-27387b6dd71b")
	if err != nil {
		return nil, fmt.Errorf("failed to parse creatorID: %w", err)
	}
	now := time.Now()

	tagsToInsert := make([]tagInsertData, len(tags))
	createdTags := make([]CreatedTagInfo, len(tags))

	for i, tag := range tags {
		tagID := uuid.New()
		tagsToInsert[i] = tagInsertData{
			ID:        tagID,
			Name:      tag.Name,
			CreatorID: creatorID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		createdTags[i] = CreatedTagInfo{
			ID:   tagID,
			Name: tag.Name,
		}
	}

	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO tags (id, name, creator_id, created_at, updated_at)
		VALUES (:id, :name, :creator_id, :created_at, :updated_at)
	`, tagsToInsert)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert tags: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdTags, nil
}

func (r *Repository) BulkAddStampMeta(ctx context.Context, additions []StampMetaAddition) error {
	if len(additions) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	creatorID, err := uuid.Parse("3b261ff3-f940-4e2c-a626-27387b6dd71b")
	if err != nil {
		return fmt.Errorf("failed to parse creatorID: %w", err)
	}
	now := time.Now()

	var linksToInsert []stampTagLinkData
	var descriptionsToInsert []stampDescriptionData

	for _, addition := range additions {
		for _, tagID := range addition.TagIDs {
			linksToInsert = append(linksToInsert, stampTagLinkData{
				StampID:   addition.StampID,
				TagID:     tagID,
				CreatorID: creatorID,
			})
		}

		if addition.Description != "" {
			descriptionsToInsert = append(descriptionsToInsert, stampDescriptionData{
				StampID:     addition.StampID,
				Description: addition.Description,
				CreatorID:   creatorID,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
		}
	}

	if len(linksToInsert) > 0 {
		log.Printf("Inserting %d stamp-tag links", len(linksToInsert))
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO stamp_tags (stamp_id, tag_id, creator_id)
			VALUES (:stamp_id, :tag_id, :creator_id)
		`, linksToInsert)
		if err != nil {
			return fmt.Errorf("failed to bulk insert stamp_tags: %w", err)
		}
	}

	if len(descriptionsToInsert) > 0 {
		log.Printf("Inserting %d stamp descriptions", len(descriptionsToInsert))
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO stamp_descriptions (stamp_id, description, creator_id, created_at, updated_at)
			VALUES (:stamp_id, :description, :creator_id, :created_at, :updated_at)
		`, descriptionsToInsert)
		if err != nil {
			log.Printf("Failed to insert stamp_descriptions. The data was: %+v", descriptionsToInsert)

			log.Printf("Error inserting descriptions: %v", err)

			return fmt.Errorf("failed to bulk insert/update stamp_descriptions: %w", err)
		}
	}

	return tx.Commit()
}
