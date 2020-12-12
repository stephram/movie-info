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
	"Provider": "cinemaworld",
	"Movies": [
		{
			"ID":         "cw12345678",
			"Title":      "one",
			"Poster":     "",
			"IsReliable": true
		},
		{
			"ID":         "cw24682468",
			"Title":      "two",
			"Poster":     "",
			"IsReliable": true
		},
		{
			"ID":         "cw35793579",
			"Title":      "tre",
			"Poster":     "",
			"IsReliable": true
		}
	]
}`

var fwMoviesJson = `{
	"Provider": "filmworld",
	"Movies": [
		{
			"ID":         "fw12345678",
			"Title":      "one",
			"Poster":     "",
			"IsReliable": true
		},
		{
			"ID":         "fw24682468",
			"Title":      "two",
			"Poster":     "",
			"IsReliable": true
		},
		{
			"ID":         "fw35793579",
			"Title":      "tre",
			"Poster":     "",
			"IsReliable": true
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
		assert.Equal(t, 3, len(movieList))
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
