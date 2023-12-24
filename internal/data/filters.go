package data

import (
	"fmt"
	"math"
	"strings"
)

type Filters struct {
	Page     int    `schema:"page" validate:"min=1"`
	PageSize int    `schema:"pageSize" validate:"max=1000000"`
	Sort     string `schema:"sort" validate:"oneof=id -id title -title runtime -runtime year -year"`
}

func (f Filters) getOrderBySpec() string {
	field := strings.TrimPrefix(f.Sort, "-")
	order := "ASC"
	if strings.HasPrefix(f.Sort, "-") {
		order = "DESC"
	}
	return fmt.Sprintf("%s %s", field, order)
}
func (f Filters) limit() int {
	return f.PageSize
}
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type MovieFilters struct {
	Title  string   `schema:"title"`
	Genres []string `schema:"genres"`
	Filters
}

func NewMovieFilters() MovieFilters {
	return MovieFilters{
		Title:  "",
		Genres: []string{},
		Filters: Filters{
			Page:     1,
			PageSize: 20,
			Sort:     "id",
		},
	}
}

type Metadata struct {
	CurrentPage  int `json:"currentPage"`
	PageSize     int `json:"pageSize"`
	FirstPage    int `json:"firstPage"`
	LastPage     int `json:"lastPage"`
	TotalRecords int `json:"totalRecords"`
}

func calculateMetadata(page int, pageSize int, totalRecords int) Metadata {
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
