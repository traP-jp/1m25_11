package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	TagSummary struct {
		ID   uuid.UUID `db:"id"`
		Name string    `db:"name"`
	}

	CreateTagParams struct {
		Name      string    `db:"name"`
		CreatorID uuid.UUID `db:"creator_id"`
	}

	Tag struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		CreatorID uuid.UUID `db:"creator_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

func (r *Repository) GetTags(ctx context.Context) ([]*TagSummary, error) {
	tags := []*TagSummary{}
	if err := r.db.SelectContext(ctx, &tags, "SELECT id,name FROM tags"); err != nil {
		return nil, fmt.Errorf("select tags: %w", err)
	}

	return tags, nil
}

func (r *Repository) UpdateTags(ctx context.Context, tagID uuid.UUID, name string) error {
	if _, err := r.db.ExecContext(ctx, `UPDATE tags SET name = ? WHERE id = ?`, name, tagID); err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

func (r *Repository) DeleteTags(ctx context.Context, tagID uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id=?`, tagID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (r *Repository) CreateTags(ctx context.Context, params CreateTagParams)(uuid.UUID, error){
	tagID, _ := uuid.NewV7()
	now := time.Now()
	if _, err := r.db.ExecContext(ctx, "INSERT INTO tags(id, name, creator_id, created_at, updated_at) VALUES(?,?,?,?,?)", tagID, params.Name, params.CreatorID, now, now); err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert tag: %w", err)
	}

	return tagID, nil
}

func (r *Repository) GetTagsByStampID(ctx context.Context, stampID uuid.UUID) ([]*TagSummary, error) {
	tagsummaries := []*TagSummary{}
	if err := r.db.SelectContext(ctx, &tagsummaries, "SELECT tags.id, tags.name FROM tags JOIN stamp_tags ON stamp_tags.tag_id = tags.id WHERE stamp_tags.stamp_id = ?", stampID); err != nil {
		return nil, fmt.Errorf("select tags by stampID: %w", err)
	}

	return tagsummaries, nil
}

func (r *Repository) GetTagDetilsByStampID(ctx context.Context, stampID uuid.UUID) ([]*Tag, error) {
	tag := []*Tag{}
	if err := r.db.SelectContext(ctx, &tag, "SELECT tags.id, tags.name, tags.creator_id, tags.created_at, tags.updated_at FROM tags JOIN stamp_tags ON stamp_tags.tag_id = tags.id WHERE stamp_tags.stamp_id = ?", stampID); err != nil {
		return nil, fmt.Errorf("select tag details by stampID: %w", err)
	}

	return tag, nil
}

var (
	ErrTagConflict   = errors.New("tag with this name already exists")
	ErrTagNotFound   = errors.New("tag not found")
	ErrAdminNotFound = errors.New("admin not found")
)
