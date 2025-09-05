package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	// users table
	User struct {
		ID                      uuid.UUID         `db:"-"`
		IsAdmin                 bool              `db:"-"`
		StampsUserOwned         []StampSummary    `db:"-"`
		TagsUserCreated         []TagSummary      `db:"-"`
		StampsUserTagged        []StampTagSummary `db:"-"`
		DescriptionsUserCreated []StampSummary    `db:"-"`
	}

	CreateUserParams struct {
		Name  string
		Email string
	}
)

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	user.ID = userID
	user.IsAdmin = false // 仮でfalseに設定
	// isAdmin, err := r.IsAdmin(ctx, userID)
	// if err != nil {
	// 	return nil, fmt.Errorf("check admin: %w", err)
	// }
	// user.IsAdmin = isAdmin.(bool)

	if err := r.db.SelectContext(ctx, &user.StampsUserOwned, "SELECT id , name ,file_id FROM stamps WHERE creator_id = ?", userID); err != nil {
		return nil, fmt.Errorf("select stamps by creatorID: %w", err)
	}
	if err := r.db.SelectContext(ctx, &user.TagsUserCreated, "SELECT id, name FROM tags WHERE creator_id = ?", userID); err != nil {
		return nil, fmt.Errorf("select tags by creatorID: %w", err)
	}
	if err := r.db.SelectContext(ctx, &user.StampsUserTagged, `SELECT 
			s.id AS "stamp.id", s.name AS "stamp.name", s.file_id AS "stamp.file_id",
			t.id AS "tag.id", t.name AS "tag.name"
			FROM stamp_tags AS st
			JOIN stamps AS s ON st.stamp_id = s.id
			JOIN tags AS t ON st.tag_id = t.id
			WHERE st.creator_id = ?`, userID); err != nil {
		return nil, fmt.Errorf("select stamp_tags by creatorID: %w", err)
	}
	if err := r.db.SelectContext(ctx, &user.DescriptionsUserCreated, `SELECT
			s.id , s.name , s.file_id 
			FROM stamp_descriptions AS d
			JOIN stamps AS s ON d.stamp_id = s.id 
			WHERE d.creator_id = ?`, userID); err != nil {
		return nil, fmt.Errorf("select stamp_descriptions by creatorID: %w", err)
	}

	return user, nil
}
