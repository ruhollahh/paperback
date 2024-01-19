package service

import (
	"github.com/ruhollahh/paperback/pkg/errsx"
	"github.com/ruhollahh/paperback/pkg/validation"
	"math"
	"strings"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f Filters) Validate(errors errsx.Map) error {
	if f.Page < 0 {
		errors.Set("page", "must be greater than zero")
	}
	if f.Page > 10_000_000 {
		errors.Set("page", "must be a maximum of 10 million")
	}
	if f.PageSize < 0 {
		errors.Set("page_size", "must be greater than zero")
	}
	if f.PageSize > 100 {
		errors.Set("page_size", "must be a maximum of 100")
	}
	if !validation.PermittedValue(f.Sort, f.SortSafeList...) {
		errors.Set("sort", "invalid sort value")
	}

	return errors
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func NewMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
