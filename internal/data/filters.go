package data

type Filters struct {
	Page     int    `schema:"page" validate:"min=1"`
	PageSize int    `schema:"pageSize" validate:"max=1000000"`
	Sort     string `schema:"sort" validate:"oneof=id -id title -title runtime -runtime year -year"`
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
