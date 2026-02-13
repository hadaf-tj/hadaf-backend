package filters

import (
	"fmt"
	"time"
)

type InstitutionFilter struct {
	Name          string     `form:"name"`
	Type          string     `form:"type"`
	City          string     `form:"city"`
	Region        string     `form:"region"`
	Address       string     `form:"address"`
	Phone         string     `form:"phone"`
	Email         string     `form:"email"`
	ActivityHours string     `form:"activity_hours"`
	CreatedAtFrom *time.Time `form:"created_at_from"`
	CreatedAtTo   *time.Time `form:"created_at_to"`
	NeedsCount    *int       `form:"needs_count"`
	OrderBy       string     `form:"order_by"`
	IsDeleted     bool       `form:"is_deleted"`
	DeletedAtFrom *time.Time `form:"deleted_at_from"`
	DeletedAtTo   *time.Time `form:"deleted_at_to"`
	Lat           *float64   `form:"lat"`
	Lng           *float64   `form:"lng"`
}

func BuildGetAllInstitutionFilter(filter InstitutionFilter) (string, []interface{}) {
	filterQuery := ""
	var args []interface{}

	// Индекс параметра чтобы не было путаницы с последовательностью аргументов
	idx := 1

	// is_deleted — всегда первый
	filterQuery += fmt.Sprintf(" WHERE is_deleted = $%d", idx)
	args = append(args, filter.IsDeleted)
	idx++

	if filter.Name != "" {
		filterQuery += fmt.Sprintf(" AND name ILIKE $%d", idx)
		args = append(args, "%"+filter.Name+"%")
		idx++
	}

	if filter.Type != "" {
		filterQuery += fmt.Sprintf(" AND type ILIKE $%d", idx)
		args = append(args, "%"+filter.Type+"%")
		idx++
	}

	if filter.City != "" {
		filterQuery += fmt.Sprintf(" AND city ILIKE $%d", idx)
		args = append(args, "%"+filter.City+"%")
		idx++
	}

	if filter.Region != "" {
		filterQuery += fmt.Sprintf(" AND region ILIKE $%d", idx)
		args = append(args, "%"+filter.Region+"%")
		idx++
	}

	if filter.Address != "" {
		filterQuery += fmt.Sprintf(" AND address ILIKE $%d", idx)
		args = append(args, "%"+filter.Address+"%")
		idx++
	}

	if filter.Phone != "" {
		filterQuery += fmt.Sprintf(" AND phone ILIKE $%d", idx)
		args = append(args, "%"+filter.Phone+"%")
		idx++
	}

	if filter.Email != "" {
		filterQuery += fmt.Sprintf(" AND email ILIKE $%d", idx)
		args = append(args, "%"+filter.Email+"%")
		idx++
	}

	if filter.ActivityHours != "" {
		filterQuery += fmt.Sprintf(" AND activity_hours ILIKE $%d", idx)
		args = append(args, "%"+filter.ActivityHours+"%")
		idx++
	}

	// ORDER BY — ТОЛЬКО whitelist
	switch filter.OrderBy {
	case "asc":
		filterQuery += " ORDER BY id ASC"
	case "desc":
		filterQuery += " ORDER BY id DESC"
	}

	return filterQuery, args
}
