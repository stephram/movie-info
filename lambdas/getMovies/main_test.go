package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"movie-info/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}

func TestPrivates(t *testing.T) {
	movieItems := []*models.MovieItem{
		{
			ID:         "cw12345678",
			Title:      "one",
			Poster:     "",
			IsReliable: true,
		},
		{
			ID:         "cw24682468",
			Title:      "two",
			Poster:     "",
			IsReliable: true,
		},
		{
			ID:         "cw35793579",
			Title:      "tre",
			Poster:     "",
			IsReliable: true,
		},
	}
	t.Run("Set Reliable", func(t *testing.T) {
		for _, movieItem := range movieItems {
			assert.True(t, movieItem.IsReliable)
		}

		_movieItems := setReliable(false, movieItems)

		assert.NotNil(t, _movieItems)

		for _, movieItem := range _movieItems {
			assert.False(t, movieItem.IsReliable)
		}

		for _, movieItem := range movieItems {
			assert.False(t, movieItem.IsReliable)
		}
	})

	t.Run("updateMovieIDs", func(t *testing.T) {
		for _, movieItem := range movieItems {
			assert.Equal(t, "", movieItem.MovieID)
		}

		movieMap := make(map[string]*models.MovieItem)
		aggregate(movieItems, movieMap)
		updateMovieIDs(movieMap)

		for _, v := range movieMap {
			assert.Equal(t, "", v.ID)
			assert.NotEqual(t, "", v.MovieID)
		}
	})
}
