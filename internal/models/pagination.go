package models

// InstitutionPage — страница списка учреждений (GET /institutions).
type InstitutionPage struct {
	Items  []*Institution `json:"items"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// EventPage — страница списка событий (GET /events).
type EventPage struct {
	Items  []*EventResponse `json:"items"`
	Total  int64            `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}
