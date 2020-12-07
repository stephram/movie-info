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
		ID:     dbMovieItem.ID,
		Title:  dbMovieItem.Title,
		Type:   dbMovieItem.Type,
		Poster: dbMovieItem.Poster,
	}
}
