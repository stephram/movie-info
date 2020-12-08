package models

type MoviesResponse struct {
	Provider string
	Movies   []*MovieItem
}

type MovieInfoResponse struct {
	ID     string
	Title  string
	Poster string
	Price  float64
}

type MovieItem struct {
	ID         string
	Title      string
	Type       string
	Poster     string
	Provider   string
	MovieID    string
	Price      float64
	IsReliable bool
}

func ConvertToMovieItem(movieInfoResponse MovieInfoResponse) *MovieItem {
	return &MovieItem{
		ID:     movieInfoResponse.ID,
		Title:  movieInfoResponse.Title,
		Poster: movieInfoResponse.Poster,
		Price:  movieInfoResponse.Price,
	}
}

func SetReliable(isReliable bool, movieItems []*MovieItem) []*MovieItem {
	for _, movieItem := range movieItems {
		movieItem.IsReliable = isReliable
	}
	return movieItems
}
