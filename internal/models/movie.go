package models

type MoviesResponse struct {
	Provider string
	Movies   []MovieItem
}

type MovieItem struct {
	ID         string
	Title      string
	Type       string
	Poster     string
	IsReliable bool `json:",omitempty"`
}