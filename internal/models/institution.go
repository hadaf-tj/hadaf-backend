package models

import "time"

type Institution struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	City          string     `json:"city"`
	Region        string     `json:"region"`
	Address       string     `json:"address"`
	Phone         string     `json:"phone"`
	Email         string     `json:"email"`
	Description   string     `json:"description"`
	ActivityHours string     `json:"activity_hours"`
	Latitude      *float64   `json:"latitude"`  // Указатель, т.к. координаты могут быть не заданы
	Longitude     *float64   `json:"longitude"` // Указатель
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}