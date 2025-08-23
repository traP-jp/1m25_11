package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type(
	ResponseStamp struct {
		ID           uuid.UUID    `json:"id"`
   		Name         string    `json:"name"`
  		CreatorID    uuid.UUID   `json:"creatorId"`
   		FileID       uuid.UUID   `json:"fileId"`
   		IsUnicode    bool      `json:"isUnicode"`
    	CreatedAt    time.Time `json:"createdAt"`
    	UpdatedAt    time.Time `json:"updatedAt"`
    	HasThumbnail bool      `json:"hasThumbnail"`
	}

)


func (r *Repository) SaveStamp(ctx context.Context, stamps []*ResponseStamp) error {

	existingIDs := make(map[uuid.UUID]time.Time)
	
	for _, resStamp := range stamps {
		
		id := resStamp.ID
		updatedAt, err := r.FindByID(ctx, id); 
		if err == sql.ErrNoRows {
			err := r.InsertStamp(ctx, resStamp)	
			if err != nil{
				log.Printf("failed to insert stamp:%v", err)

				return fmt.Errorf("failed to insert stamp:%w",err)
			}
		}else if err != nil{
			log.Printf("failed to find stamp:%v", err)

			return fmt.Errorf("failed to find stamp:%w",err)
		}else{
			if updatedAt.Equal(resStamp.UpdatedAt) {
			err := r.UpdateStamp(ctx, resStamp)
			if err != nil {
				log.Printf("failed to update stamp:%v", err)
				
				return fmt.Errorf("failed to update stamp:%w",err)
			}
		}

		}

		
		return nil
        
	}

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID)(time.Time, error) {
	updatedAt := time.Time{}
	err := r.db.GetContext(ctx, &updatedAt,"SELECT updated_at FROM stamps WHERE id = ?",id); 	

	 if err != nil {      
		return time.Time{}, err
    }
    
    return updatedAt,nil
}



func (r *Repository) InsertStamp(ctx context.Context, resStamp *ResponseStamp) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO stamps(id,name,creator_id,file_id,is_unicode,created_at,updated_at,count_monthly,count_total) VALUES(?,?,?,?,?,?,?,0,0)", resStamp.ID, resStamp.Name, resStamp.CreatorID, resStamp.FileID, resStamp.IsUnicode, resStamp.CreatedAt, resStamp.UpdatedAt); err != nil {
		return err
	}

	return nil
	
}

func (r *Repository) UpdateStamp(ctx context.Context, resStamp *ResponseStamp) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE stamps SET name = ?,creator_id = ?,file_id = ?,is_unicode = ?,updated_at = ? WHERE id = ?", resStamp.Name, resStamp.CreatorID, resStamp.FileID, resStamp.IsUnicode, time.Now(), resStamp.ID); err != nil {
		return err
	}

	return nil
}