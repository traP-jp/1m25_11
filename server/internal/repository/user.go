package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type (
	// users table
	User struct {
		ID                      uuid.UUID `db:"-"`
		IsAdmin                 bool                `db:"-"`
		StampsUserOwned         []StampSummary      `db:"-"`
		TagsUserCreated         []TagSummary        `db:"-"`
		StampsUserTagged        []StampTagSummary   `db:"-"`
		DescriptionsUserCreated []DescriptionSummary `db:"-"`
	}

	CreateUserParams struct {
		Name  string
		Email string
	}
)

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	user.ID = userID
	log.Print(0)
	user.IsAdmin = false // 仮でfalseに設定
	// isAdmin, err := r.IsAdmin(ctx, userID)
	// if err != nil {
	// 	return nil, fmt.Errorf("check admin: %w", err)
	// }
	// user.IsAdmin = isAdmin.(bool)

	if err := r.db.SelectContext(ctx, &user.StampsUserOwned, "SELECT id AS stamp_id, name AS stamp_name, file_id FROM stamps WHERE creator_id = ?", userID); err != nil {
		return nil, fmt.Errorf("select stamps by creatorID: %w", err)
	}
	log.Print(1)
	if err := r.db.SelectContext(ctx, &user.TagsUserCreated, "SELECT id AS tag_id, name AS tag_name, FROM tags WHERE creator_id = ?", userID); err != nil {
		return nil, fmt.Errorf("select tags by creatorID: %w", err)
	}
	log.Print(2)
	if err := r.db.SelectContext(ctx, &user.StampsUserTagged, `SELECT 
			s.id AS "stamp_id", s.name AS "stamp_name", s.file_id AS "file_id",
			t.id AS "tag_id", t.name AS "tag_name"
			FROM stamp_tags AS st
			JOIN stamps AS s ON st.stamp_id = s.id
			JOIN tags AS t ON st.tag_id = t.id
			WHERE st.creator_id = ?`, userID); err != nil {
		return nil, fmt.Errorf("select stamp_tags by creatorID: %w", err)
	}
	log.Print(3)
	if err := r.db.SelectContext(ctx, &user.DescriptionsUserCreated, `SELECT
			s.id AS "stamp_id", s.name AS "stamp_name", s.file_id AS "file_id",
			d.id AS "description_id"
			FROM stamp_descriptions AS d
			JOIN stamps AS s ON d.stamp_id = s.id 
			WHERE d.creator_id = ?`, userID); err != nil {
		return nil, fmt.Errorf("select stamp_descriptions by creatorID: %w", err)
	}
	log.Print(4)	

	return user, nil
}
