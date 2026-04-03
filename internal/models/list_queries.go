package models

// InstitutionListQuery — параметры списка учреждений (Handler → Service → Repository).
type InstitutionListQuery struct {
	Search  string
	Type    string
	UserLat float64
	UserLng float64
	SortBy  string
	Limit   int
	Offset  int
}

// EventListQuery — параметры списка событий.
type EventListQuery struct {
	UserID int
	Limit  int
	Offset int
}

// EventDetailQuery — параметры карточки события (GET /events/:id). ViewerUserID = 0 без авторизации.
type EventDetailQuery struct {
	EventID      int
	ViewerUserID int
}
