package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	// users table
	User struct {
		ID               uuid.UUID `db:"id"`
		IsAdmin          bool
		StampsUserOwned  []StampSummary
		TagsUserCreated  []TagSummary
		StampsUserTagged []struct {
			Stamp StampSummary
			Tag   TagSummary
		}
		DescriptionsUserCreated []struct {
			Stamp         StampSummary
			DescriptionID uuid.UUID
		}
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
	isAdmin := false // 仮でfalseに設定
	stampsUserOwned, err := r.getStampsByCreatorID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get stamps user owned: %w", err)
	}
	tagsUserCreated, err := r.getTagsByCreatorID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get tags user created: %w", err)
	}
	stampsUserTagged, err := r.getStampTagsByCreatorID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get stamps user tagged: %w", err)
	}
	descriptionsUserCreated, err := r.getDescriptionsByCreatorID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get descriptions user created: %w", err)
	}

	return user, nil
}
