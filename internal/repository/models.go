package repository

import "movie-info/internal/models"

type DbMovieItem struct {
	ID       string
	MovieID  string
	Provider string
	Title    string
	Type     string
	Poster   string
}

func convertToDbMovieItem(movieProvider string, movieItem models.MovieItem) *DbMovieItem {
	return &DbMovieItem{
		ID:       movieItem.ID,
		MovieID:  movieItem.ID[2:],
		Provider: movieProvider,
		Title:    movieItem.Title,
		Type:     movieItem.Type,
		Poster:   movieItem.Poster,
	}
}

func convertToMovieItem(dbMovieItem DbMovieItem) *models.MovieItem {
	return &models.MovieItem{
		ID:      dbMovieItem.ID,
		Title:   dbMovieItem.Title,
		Type:    dbMovieItem.Type,
		Poster:  dbMovieItem.Poster,
		MovieID: dbMovieItem.MovieID,
	}
}

func convertToMovieItems(dbMovieItems []DbMovieItem) []*models.MovieItem {
	movieItems := make([]*models.MovieItem, 0)

	for _, dbMovieItem := range dbMovieItems {
		movieItems = append(movieItems, convertToMovieItem(dbMovieItem))
	}
	return movieItems
}
