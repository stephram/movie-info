package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"movie-info/internal/models"
	"testing"
)

var cwMovieJson = `{
  "ID": "cw2488496",
  "Title": "Star Wars: Episode VII - The Force Awakens",
  "Type": "movie",
  "Poster": "https://m.media-amazon.com/images/M/MV5BOTAzODEzNDAzMl5BMl5BanBnXkFtZTgwMDU1MTgzNzE@._V1_SX300.jpg",
  "Price": 25.50
}`

var fwMovieJson = `{
  "ID": "fw2488496",
  "Title": "Star Wars: Episode VII - The Force Awakens",
  "Type": "movie",
  "Poster": "https://m.media-amazon.com/images/M/MV5BOTAzODEzNDAzMl5BMl5BanBnXkFtZTgwMDU1MTgzNzE@._V1_SX300.jpg",
  "Price": 19.50
}`

func TestHandler(t *testing.T) {
	t.Run("JSON decode OK", func(t *testing.T) {
		var jErr error

		cwMovieItem := &models.MovieItem{}
		jErr = json.Unmarshal([]byte(cwMovieJson), &cwMovieItem)
		assert.Nil(t, jErr)

		fwMovieItem := &models.MovieItem{}
		jErr = json.Unmarshal([]byte(fwMovieJson), &fwMovieItem)
		assert.Nil(t, jErr)

		assert.Equal(t, "cw2488496", cwMovieItem.ID)
		assert.Equal(t, "fw2488496", fwMovieItem.ID)
		assert.Equal(t, cwMovieItem.Title, fwMovieItem.Title)
	})
}
