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
			Poster:     "one-poster",
			IsReliable: true,
		},
		{
			ID:         "fw24682468",
			Title:      "two",
			Poster:     "two-poster",
			IsReliable: true,
		},
		{
			ID:         "fw35793579",
			Title:      "tre",
			Poster:     "tre-poster",
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
		{
			ID:       "fw24682468",
			MovieID:  "24682468",
			Provider: "Filmworld",
			Title:    "two",
			Type:     "movie",
			Poster:   "two-poster",
		},
		{
			ID:       "fw35793579",
			MovieID:  "35793579",
			Provider: "Filmworld",
			Title:    "tre",
			Type:     "movie",
			Poster:   "tre-poster",
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

	t.Run("convertToMovieItems", func(t *testing.T) {
		movieItems := convertToMovieItems(dbMovieItems)

		assert.Equal(t, len(dbMovieItems), len(movieItems))

		for i, dbMovieItem := range dbMovieItems {
			assert.Equal(t, dbMovieItem.ID, movieItems[i].ID)
			assert.Equal(t, dbMovieItem.Title, movieItems[i].Title)
			assert.Equal(t, dbMovieItem.Poster, movieItems[i].Poster)
		}
	})
}
