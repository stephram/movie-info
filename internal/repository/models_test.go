package repository

import (
	"github.com/stretchr/testify/assert"
	"movie-info/internal/models"
	"testing"
)

func TestConvertToDbMovieItem(t *testing.T) {
	movieItems := []models.MovieItem{
		{
			ID:         "fw12345678",
			Title:      "one",
			Poster:     "",
			IsReliable: true,
		},
		{
			ID:         "fw24682468",
			Title:      "two",
			Poster:     "",
			IsReliable: true,
		},
		{
			ID:         "fw35793579",
			Title:      "tre",
			Poster:     "",
			IsReliable: true,
		},
	}

	dbMovieItems := []DbMovieItem{
		{
			ID:       "fw12345678",
			MovieID:  "12345678",
			Provider: "Filmworld",
			Title:    "one",
			Type:     "movie",
			Poster:   "one-poster",
		},
	}

	t.Run("convertToDbMovieItem", func(t *testing.T) {
		movieItem := movieItems[0]

		dbMovieItem := convertToDbMovieItem("Filmworld", movieItem)
		assert.NotNil(t, dbMovieItem)
		assert.Equal(t, movieItem.ID, dbMovieItem.ID)
		assert.Equal(t, movieItem.ID[2:], dbMovieItem.MovieID)
		assert.Equal(t, movieItem.Title, dbMovieItem.Title)
		assert.Equal(t, movieItem.Poster, dbMovieItem.Poster)
	})

	t.Run("convertToMovieItem", func(t *testing.T) {
		dbMovieItem := dbMovieItems[0]

		movieItem := convertToMovieItem(dbMovieItem)
		assert.NotNil(t, movieItem)
		assert.Equal(t, dbMovieItem.ID, movieItem.ID)
		assert.Equal(t, dbMovieItem.Title, movieItem.Title)
	})
}
