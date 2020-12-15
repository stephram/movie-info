package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "log"
	"movie-info/internal/models"
	"movie-info/internal/repository"
	"movie-info/internal/utils"
	"net/http"
	"os"
)

var (
	// Should be reading these from SSM
	movieDataEndpoint = os.Getenv("MOVIE_DATA_ENDPOINT")
	movieDataApiKey   = os.Getenv("MOVIE_DATA_API_KEY")
	// movieProviderNames = os.Getenv("MOVIE_PROVIDERS")
	movieTable = os.Getenv("MOVIE_TABLE")

	log = utils.GetLogger()

	client *http.Client
	repo   repository.Repository
)

func init() {
	client = &http.Client{}
	repo = repository.New()
}

// NEEDS REFACTORING. LOOK AWAY.

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info().Msgf("RequestID=%s, RequestTime=%s, Path=%s, QueryStringParameters=%+v",
		request.RequestContext.RequestID,
		request.RequestContext.RequestTime,
		request.Path,
		request.QueryStringParameters)

	movieID := request.QueryStringParameters["movieId"]

	movieItems, rErr := repo.GetMoviesByMovieID(movieTable, movieID)
	if rErr != nil {
		log.Err(rErr).Msgf("failed to read movieID %s from repository", movieID)
	}
	log.Info().Interface("movieItems", movieItems).Msgf("read movies for movieID %s", movieID)

	movieList := make([]*models.MovieItem, 0)

	for _, movieItem := range movieItems {
		uri := movieDataEndpoint + "/" + movieItem.Provider + "/movie/" + movieItem.ID

		req, _ := http.NewRequest("GET", uri, nil)
		req.Header.Add("x-api-key", movieDataApiKey)

		mRes, mErr := client.Do(req)
		if mErr != nil {
			return utils.CreateApiGwResponse(500, mErr.Error()), nil
		}
		defer mRes.Body.Close()

		if mRes.StatusCode != 200 {
			log.Error().Msgf("Status %d (%s) : failed to retrieve movie from %s", mRes.StatusCode, mRes.Status, uri)
			movieList = append(movieList, updateMovieItem(false, movieID, movieItem))
			continue
		}
		log.Info().Msgf("retrieved movie info from %s", uri)

		var movieInfoResponse models.MovieInfoResponse
		jErr := json.NewDecoder(mRes.Body).Decode(&movieInfoResponse)
		if jErr != nil {
			log.Err(jErr).Msgf("failed to decode MovieInfoResponse from url %s", uri)
			movieList = append(movieList, updateMovieItem(false, movieID, movieItem))
			continue
		}
		_movieItem := models.ConvertToMovieItem(movieInfoResponse)
		_movieItem.Provider = movieItem.Provider
		movieList = append(movieList, updateMovieItem(true, movieID, _movieItem))

		rErr := repo.UpdateMovieItem(movieTable, _movieItem)
		if rErr != nil {
			log.Err(rErr).Msgf("failed to update MovieItem")
		}
	}
	payload, jsonErr := json.Marshal(movieList)
	if jsonErr != nil {
		return utils.CreateApiGwResponse(500, jsonErr.Error()), nil
	}
	log.Info().Interface("movies", movieList).Msgf("Success reading movies for MovieID %s", movieID)

	// OK, return the Movies
	return utils.CreateApiGwResponse(200, string(payload)), nil
}

func updateMovieItem(isReliable bool, movieID string, movieItem *models.MovieItem) *models.MovieItem {
	movieItem.IsReliable = isReliable
	movieItem.MovieID = movieID
	return movieItem
}

func main() {
	lambda.Start(handler)
}
