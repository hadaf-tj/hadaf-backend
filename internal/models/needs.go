package models

import "time"

type Need struct {
	ID            int       `json:"id" db:"id"`
	InstitutionID int       `json:"institution_id" db:"institution_id"`
	CategoryID    *int      `json:"category_id" db:"category_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Unit          string    `json:"unit" db:"unit"`
	RequiredQty   float64   `json:"required_qty" db:"required_qty"`
	ReceivedQty   float64   `json:"received_qty" db:"received_qty"`
	Urgency       string    `json:"urgency" db:"urgency"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}