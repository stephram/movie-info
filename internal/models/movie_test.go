package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrivates(t *testing.T) {
	movieItems := []*MovieItem{
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

		_movieItems := SetReliable(false, movieItems)

		assert.NotNil(t, _movieItems)

		for _, movieItem := range _movieItems {
			assert.False(t, movieItem.IsReliable)
		}

		for _, movieItem := range movieItems {
			assert.False(t, movieItem.IsReliable)
		}
	})
}
