package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"movie-info/internal/models"
	"movie-info/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var cwMoviesJson = `{
  "Provider": "Cinema World",
  "Movies": [
    {
      "ID": "cw2488496",
      "Title": "Star Wars: Episode VII - The Force Awakens",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOTAzODEzNDAzMl5BMl5BanBnXkFtZTgwMDU1MTgzNzE@._V1_SX300.jpg"
    },
    {
      "ID": "cw2527336",
      "Title": "Star Wars: Episode VIII - The Last Jedi",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMjQ1MzcxNjg4N15BMl5BanBnXkFtZTgwNzgwMjY4MzI@._V1_SX300.jpg"
    },
    {
      "ID": "cw2527338",
      "Title": "Star Wars: Episode IX - The Rise of Skywalker",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMDljNTQ5ODItZmQwMy00M2ExLTljOTQtZTVjNGE2NTg0NGIxXkEyXkFqcGdeQXVyODkzNTgxMDg@._V1_SX300.jpg"
    },
    {
      "ID": "cw3748528",
      "Title": "Rogue One: A Star Wars Story",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMjEwMzMxODIzOV5BMl5BanBnXkFtZTgwNzg3OTAzMDI@._V1_SX300.jpg"
    },
    {
      "ID": "cw3778644",
      "Title": "Solo: A Star Wars Story",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOTM2NTI3NTc3Nl5BMl5BanBnXkFtZTgwNzM1OTQyNTM@._V1_SX300.jpg"
    },
    {
      "ID": "cw0076759",
      "Title": "Star Wars: Episode IV - A New Hope",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNzVlY2MwMjktM2E4OS00Y2Y3LWE3ZjctYzhkZGM3YzA1ZWM2XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg"
    },
    {
      "ID": "cw0080684",
      "Title": "Star Wars: Episode V - The Empire Strikes Back",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYmU1NDRjNDgtMzhiMi00NjZmLTg5NGItZDNiZjU5NTU4OTE0XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg"
    },
    {
      "ID": "cw0086190",
      "Title": "Star Wars: Episode VI - Return of the Jedi",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOWZlMjFiYzgtMTUzNC00Y2IzLTk1NTMtZmNhMTczNTk0ODk1XkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg"
    },
    {
      "ID": "cw0120915",
      "Title": "Star Wars: Episode I - The Phantom Menace",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYTRhNjcwNWQtMGJmMi00NmQyLWE2YzItODVmMTdjNWI0ZDA2XkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg"
    },
    {
      "ID": "cw0121765",
      "Title": "Star Wars: Episode II - Attack of the Clones",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMDAzM2M0Y2UtZjRmZi00MzVlLTg4MjEtOTE3NzU5ZDVlMTU5XkEyXkFqcGdeQXVyNDUyOTg3Njg@._V1_SX300.jpg"
    },
    {
      "ID": "cw0121766",
      "Title": "Star Wars: Episode III - Revenge of the Sith",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNTc4MTc3NTQ5OF5BMl5BanBnXkFtZTcwOTg0NjI4NA@@._V1_SX300.jpg"
    }
  ]
}`

var fwMoviesJson = `{
  "Provider": "Film World",
  "Movies": [
    {
      "ID": "fw2488496",
      "Title": "Star Wars: Episode VII - The Force Awakens",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOTAzODEzNDAzMl5BMl5BanBnXkFtZTgwMDU1MTgzNzE@._V1_SX300.jpg"
    },
    {
      "ID": "fw2527336",
      "Title": "Star Wars: Episode VIII - The Last Jedi",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMjQ1MzcxNjg4N15BMl5BanBnXkFtZTgwNzgwMjY4MzI@._V1_SX300.jpg"
    },
    {
      "ID": "fw2527338",
      "Title": "Star Wars: Episode IX - The Rise of Skywalker",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMDljNTQ5ODItZmQwMy00M2ExLTljOTQtZTVjNGE2NTg0NGIxXkEyXkFqcGdeQXVyODkzNTgxMDg@._V1_SX300.jpg"
    },
    {
      "ID": "fw3748528",
      "Title": "Rogue One: A Star Wars Story",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMjEwMzMxODIzOV5BMl5BanBnXkFtZTgwNzg3OTAzMDI@._V1_SX300.jpg"
    },
    {
      "ID": "fw3778644",
      "Title": "Solo: A Star Wars Story",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOTM2NTI3NTc3Nl5BMl5BanBnXkFtZTgwNzM1OTQyNTM@._V1_SX300.jpg"
    },
    {
      "ID": "fw0076759",
      "Title": "Star Wars: Episode IV - A New Hope",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNzVlY2MwMjktM2E4OS00Y2Y3LWE3ZjctYzhkZGM3YzA1ZWM2XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg"
    },
    {
      "ID": "fw0080684",
      "Title": "Star Wars: Episode V - The Empire Strikes Back",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYmU1NDRjNDgtMzhiMi00NjZmLTg5NGItZDNiZjU5NTU4OTE0XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg"
    },
    {
      "ID": "fw0086190",
      "Title": "Star Wars: Episode VI - Return of the Jedi",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOWZlMjFiYzgtMTUzNC00Y2IzLTk1NTMtZmNhMTczNTk0ODk1XkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg"
    },
    {
      "ID": "fw0120915",
      "Title": "Star Wars: Episode I - The Phantom Menace",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYTRhNjcwNWQtMGJmMi00NmQyLWE2YzItODVmMTdjNWI0ZDA2XkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg"
    },
    {
      "ID": "fw0121765",
      "Title": "Star Wars: Episode II - Attack of the Clones",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMDAzM2M0Y2UtZjRmZi00MzVlLTg4MjEtOTE3NzU5ZDVlMTU5XkEyXkFqcGdeQXVyNDUyOTg3Njg@._V1_SX300.jpg"
    },
    {
      "ID": "fw0121766",
      "Title": "Star Wars: Episode III - Revenge of the Sith",
      "Type": "movie",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNTc4MTc3NTQ5OF5BMl5BanBnXkFtZTcwOTg0NjI4NA@@._V1_SX300.jpg"
    }
  ]
}`

func init() {
	os.Setenv("MOVIE_DATA_ENDPOINT", "http://localhost:8080/api")
	os.Setenv("MOVIE_DATA_API_KEY", "1234567876543210")
	os.Setenv("MOVIE_PROVIDERS", "[\"cinemaworld\", \"filmworld\"]")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func Test_handler(t *testing.T) {
	mockRepo := &mocks.Repository{}
	repo = mockRepo

	t.Run("handler_TestWithMock", func(t *testing.T) {
		mockRepo.On("UpdateProviderMovies", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			if strings.Contains(r.URL.Path, "cinemaworld") {
				fmt.Fprint(w, cwMoviesJson)
			} else if strings.Contains(r.URL.Path, "filmworld") {
				fmt.Fprint(w, fwMoviesJson)
			} else {
				fmt.Fprint(w, "{}")
			}

		}))
		defer ts.Close()

		movieResponse := &models.MoviesResponse{}
		jErr := json.Unmarshal([]byte(cwMoviesJson), &movieResponse)
		if jErr != nil {
			assert.FailNow(t, "failed to Unmarshal MovieResponse : %s", jErr.Error())
		}
		assert.NotNil(t, movieResponse)

		os.Setenv("MOVIE_DATA_ENDPOINT", ts.URL)

		apiGwRes, apiGwErr := handler(events.APIGatewayProxyRequest{})
		if apiGwErr != nil {
			t.Error(t, "failed with error: %s", apiGwErr)
		}

		assert.Equal(t, 200, apiGwRes.StatusCode)

		movieList := make([]*models.MovieItem, 0)
		rErr := json.Unmarshal([]byte(apiGwRes.Body), &movieList)
		assert.Nil(t, rErr)
		assert.Equal(t, 11, len(movieList))
	})

	t.Run("handler_FailDueToProviderConfiguration", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{}

		apiGwRes, apiGwErr := handler(request)
		assert.Nil(t, apiGwErr)
		assert.NotNil(t, apiGwRes)
		assert.Equal(t, 500, apiGwRes.StatusCode)
		assert.True(t, len(apiGwRes.Body) > 0)
	})
}

func Test_getProviderConfiguration(t *testing.T) {
	t.Run("getProviderConfiguration_OK", func(t *testing.T) {
		os.Setenv("MOVIE_PROVIDERS", "[\"cinemaworld\", \"filmworld\"]")

		movieProviders, err := getProviderConfiguration()
		assert.Nil(t, err)
		assert.NotNil(t, movieProviders)
		assert.Equal(t, "cinemaworld", movieProviders[0])
		assert.Equal(t, "filmworld", movieProviders[1])
	})

	t.Run("getProviderConfiguration_ERR", func(t *testing.T) {
		os.Setenv("MOVIE_PROVIDERS", "")

		movieProviders, err := getProviderConfiguration()
		assert.NotNil(t, err)
		assert.Nil(t, movieProviders)
	})
}

func Test_getProviderMovies(t *testing.T) {

}
