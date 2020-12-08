package repository

import "movie-info/internal/models"

type DbMovieItem struct {
	ID       string
	MovieID  string
	Provider string
	Title    string
	Type     string
	Poster   string
	Price    float64
}

func convertToDbMovieItem(movieProvider string, movieItem models.MovieItem) *DbMovieItem {
	return &DbMovieItem{
		ID:       movieItem.ID,
		MovieID:  movieItem.ID[2:],
		Provider: movieProvider,
		Title:    movieItem.Title,
		Type:     movieItem.Type,
		Poster:   movieItem.Poster,
		Price:    movieItem.Price,
	}
}

func convertToMovieItem(dbMovieItem DbMovieItem) *models.MovieItem {
	return &models.MovieItem{
		ID:       dbMovieItem.ID,
		Title:    dbMovieItem.Title,
		Type:     dbMovieItem.Type,
		Poster:   dbMovieItem.Poster,
		Provider: dbMovieItem.Provider,
		MovieID:  dbMovieItem.MovieID,
		Price:    dbMovieItem.Price,
	}
}

func convertToMovieItems(dbMovieItems []DbMovieItem) []*models.MovieItem {
	movieItems := make([]*models.MovieItem, 0)

	for _, dbMovieItem := range dbMovieItems {
		movieItems = append(movieItems, convertToMovieItem(dbMovieItem))
	}
	return movieItems
}
