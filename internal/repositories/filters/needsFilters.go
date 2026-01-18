package filters

import (
	"fmt"
	"time"
)

type NeedsFilter struct {
	CategoryID  *int       `form:"category_id"`
	Name        string     `form:"name"`
	Unit        string     `form:"unit"`
	RequiredQty float64    `form:"required_qty"`
	ReceivedQty float64    `form:"received_qty"`
	Urgency     string     `form:"urgency"`
	CreatedAt   *time.Time `form:"created_at"`
	IsDeleted   bool       `form:"is_deleted"`
	OrderBy     string     `form:"order_by"`
}

func GetNeedsByInstitution(filter NeedsFilter, instituteId int) (string, []interface{}) {
	filterQuery := ""
	var args []interface{}

	// Индекс параметра чтобы не было путаницы с последовательностью аргументов
	idx := 1

	// is_deleted — всегда первый
	filterQuery += fmt.Sprintf(" WHERE is_deleted = $%d", idx)
	args = append(args, filter.IsDeleted)
	idx++

	filterQuery += fmt.Sprintf(" AND institution_id = $%d", idx)
	args = append(args, instituteId)
	idx++

	if filter.Name != "" {
		filterQuery += fmt.Sprintf(" AND name ILIKE $%d", idx)
		args = append(args, "%"+filter.Name+"%")
		idx++
	}

	if filter.CategoryID != nil {
		filterQuery += fmt.Sprintf(" AND category_id = $%d", idx)
		args = append(args, filter.CategoryID)
		idx++
	}

	if filter.Unit != "" {
		filterQuery += fmt.Sprintf(" AND unit ILIKE $%d", idx)
		args = append(args, filter.Unit)
		idx++
	}

	if filter.RequiredQty != 0 {
		filterQuery += fmt.Sprintf(" AND required_qty >= $%d", idx)
		args = append(args, filter.RequiredQty)
		idx++
	}

	if filter.ReceivedQty != 0 {
		filterQuery += fmt.Sprintf(" AND received_qty >= $%d", idx)
		args = append(args, filter.ReceivedQty)
		idx++
	}

	if filter.Urgency != "" {
		filterQuery += fmt.Sprintf(" AND urgency ILIKE $%d", idx)
		args = append(args, "%"+filter.Urgency+"%")
		idx++
	}

	if filter.CreatedAt != nil {
		filterQuery += fmt.Sprintf(" AND created_at = $%d", idx)
		args = append(args, *filter.CreatedAt)
		idx++
	}

	// ORDER BY — ТОЛЬКО whitelist
	switch filter.OrderBy {
	case "asc":
		filterQuery += " ORDER BY created_at ASC"
	case "desc":
		filterQuery += " ORDER BY created_at DESC"
	}
	return filterQuery, args
}
