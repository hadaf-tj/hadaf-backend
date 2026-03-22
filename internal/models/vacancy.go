package models

import "time"

// Vacancy модель вакансий
type Vacancy struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`       // например "Волонтёрство"
	Experience  string     `json:"experience"` // например "От 1 года"
	Workload    string     `json:"workload"`   // например "1 час в неделю"
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	IsDeleted   bool       `json:"-"`
}
