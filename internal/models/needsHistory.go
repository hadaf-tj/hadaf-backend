package models

import "time"

type NeedsHistory struct {
	ID        int        `json:"id"`
	NeedID    int        `json:"need_id"`
	Comment   *string    `json:"comment"`
	CreatedAt time.Time  `json:"created_at"`
	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
}
