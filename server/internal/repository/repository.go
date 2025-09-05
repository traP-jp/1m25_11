package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}
type TagDetails struct {
	ID        uuid.UUID     `db:"id"`
	Name      string        `db:"name"`
	CreatorID uuid.UUID     `db:"creator_id"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	Stamps    []StampForTag `db:"-"`
}

type StampForTag struct {
	ID     uuid.UUID `db:"id"`
	Name   string    `db:"name"`
	FileID uuid.UUID `db:"file_id"`
}

func (r *Repository) CreateTag(context context.Context, params CreateTagParams) (any, any) {
	panic("unimplemented")
}

func (r *Repository) GetTagDetails(ctx context.Context, tagID uuid.UUID) (*TagDetails, error) {
	var tagDetails TagDetails
	err := r.db.GetContext(ctx, &tagDetails, "SELECT id, name, creator_id, created_at, updated_at FROM tags WHERE id = ?", tagID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTagNotFound
		}

		return nil, err
	}

	var stamps []StampForTag
	err = r.db.SelectContext(ctx, &stamps, `
		SELECT s.id, s.name, s.file_id FROM stamps s
		INNER JOIN stamp_tags st ON s.id = st.stamp_id
		WHERE st.tag_id = ?`, tagID)
	if err != nil {
		return nil, err
	}

	tagDetails.Stamps = stamps

	return &tagDetails, nil
}

func (r *Repository) IsAdmin(context context.Context, userID uuid.UUID) (any, error) {
	panic("unimplemented")
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}
