package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	// users table
	User struct {
		ID                      uuid.UUID `db:"id"`
		IsAdmin                 bool
		StampsUserOwned         []StampSummary
		TagsUserCreated         []TagSummary
		StampsUserTagged        []StampTagSummary
		DescriptionsUserCreated []DescriptionSummary
	}

	CreateUserParams struct {
		Name  string
		Email string
	}
)

func (r *Repository) GetUsers(ctx context.Context) ([]*User, error) {
	users := []*User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (r *Repository) CreateUser(ctx context.Context, params CreateUserParams) (uuid.UUID, error) {
	userID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate uuid: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users (id, name, email) VALUES (?, ?, ?)", userID, params.Name, params.Email); err != nil {
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	user.ID = userID
	user.isAdmin := false // 仮でfalseに設定

	if err := r.db.SelectContext(ctx, &user.StampsUserOwned, "SELECT id, name, file_id FROM stamps WHERE creator_id = ?", userID); err != nil {
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
	descriptionsUserCreated := []*DescriptionSummary{}
	if err := r.db.SelectContext(ctx, &descriptionsUserCreated, `SELECT
			s.id AS "stamp.id", s.name AS "stamp.name", s.file_id AS "stamp.file_id",
			d.id AS "description_id"
			FROM stamp_descriptions AS d
			JOIN stamps AS s ON d.stamp_id = s.id 
			WHERE d.creator_id = ?`, userID); err != nil {
		return nil, fmt.Errorf("select stamp_descriptions by creatorID: %w", err)
	}

	return user, nil
}
