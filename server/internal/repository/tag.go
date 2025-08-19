package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	// stamps table
	Tag struct {
		ID    uuid.UUID `db:"id"`
		Name  string    `db:"name"`		
	}

	CreateTagParams struct{
		Name string `db:"name"`
		CreatorID uuid.UUID `db:"creator_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

)

func (r *Repository) GetTags(ctx context.Context) ([]*Tag, error) {
	tags := []*Tag{}
	if err := r.db.SelectContext(ctx, &tags, "SELECT id,name FROM tags"); err != nil {
		return nil, fmt.Errorf("select tags: %w", err)
	}
	
	return tags, nil
}

func (r *Repository) UpdateTags(ctx context.Context, tagID uuid.UUID, name string)error{
	if _, err :=r.db.ExecContext(ctx,`UPDATE tags SET name = ? WHERE id = ?`, name, tagID); err != nil{
		return fmt.Errorf("failed to update tag: %w",err)
	}

	return nil
}

func (r *Repository) DeleteTags(ctx context.Context, tagID uuid.UUID)error{
	if _,err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id=?`, tagID); err != nil{
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (r *Repository) CreateTags(ctx context.Context, params CreateTagParams)(uuid.UUID, error){
	tagID := uuid.New()
	if _, err := r.db.ExecContext(ctx, "INSERT INTO tags(id, name, creator_id, created_at, updated_at) VALUES(?,?,?,?,?)", tagID, params.Name, params.CreatorID, params.CreatedAt, params.UpdatedAt);err != nil{
		return uuid.Nil, fmt.Errorf("failed to insert tag: %w", err )
	}

	return tagID, nil
}

func (r *Repository) GetTagsByStampID(ctx context.Context, stampID uuid.UUID)([]*Tag,error){
	tags := []*Tag{}
	if err := r.db.SelectContext(ctx, &tags, "SELECT tags.id, tags.name FROM tags JOIN stamp_tags ON stamp_tags.tag_id = tags.id WHERE stamp_tags.stamp_id = ?",stampID); err != nil{
		return nil, fmt.Errorf("select tags by stampID: %w",err )
	}

	return tags, nil
}


