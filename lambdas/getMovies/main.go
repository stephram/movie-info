package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	_ "log"
	"movie-info/internal/models"
	"movie-info/internal/repository"
	"movie-info/internal/utils"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// Should be reading these from SSM
	movieDataEndpoint string
	movieDataApiKey   string
	movieTable        string

	log = utils.GetLogger()

	client *http.Client
	repo   repository.Repository
)

func init() {
	client = &http.Client{}
	repo = repository.New()
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info().Msgf("RequestID=%s, RequestTime=%s", request.RequestContext.RequestID, request.RequestContext.RequestTime)

	movieDataEndpoint = os.Getenv("MOVIE_DATA_ENDPOINT")
	movieDataApiKey = os.Getenv("MOVIE_DATA_API_KEY")
	movieTable = os.Getenv("MOVIE_TABLE")

	movieProviders, pErr := getProviderConfiguration()
	if pErr != nil {
		return utils.CreateApiGwResponse(500, pErr.Error()), nil
	}

	movieMap, mErr := processProviders(movieProviders)
	if mErr != nil {
		return utils.CreateApiGwResponse(500, mErr.Error()), nil
	}

	// Prepare IDs and Flatten Map to a Slice with unique items.
	movieList, prErr := prepareResponse(movieMap)
	if prErr != nil {
		return utils.CreateApiGwResponse(500, pErr.Error()), nil
	}

	// Marshal to JSON.
	payload, jsonErr := json.Marshal(movieList)
	if jsonErr != nil {
		return utils.CreateApiGwResponse(500, jsonErr.Error()), nil
	}
	// OK, return the Movies
	return utils.CreateApiGwResponse(200, string(payload)), nil
}

func processProviders(movieProviders []string) (map[string][]*models.MovieItem, error) {
	movieMap := make(map[string][]*models.MovieItem)

	for _, movieProvider := range movieProviders {
		pErr := getProviderMovies(movieProvider, movieMap)
		if pErr != nil {
			return nil, pErr
		}

		uErr := updateCacheForProvider(movieProvider, movieMap[movieProvider])
		if uErr != nil {
			log.Err(uErr).Msgf(errors.Wrapf(uErr, "failed to update database").Error())
		}
	}
	return movieMap, nil
}

func prepareResponse(movieMap map[string][]*models.MovieItem) ([]*models.MovieItem, error) {
	movieList := make([]*models.MovieItem, 0)
	flattened := make(map[string]*models.MovieItem)

	for _, v := range movieMap {
		for _, vv := range v {
			vv.MovieID = vv.ID[2:]
			vv.ID = ""
			flattened[vv.MovieID] = vv
		}
	}
	for _, v := range flattened {
		movieList = append(movieList, v)
	}
	return movieList, nil
}

func getProviderConfiguration() ([]string, error) {
	movieProviderNames := os.Getenv("MOVIE_PROVIDERS")

	var movieProviders []string
	jsonErr := json.Unmarshal([]byte(movieProviderNames), &movieProviders)
	if jsonErr != nil {
		return nil, errors.Wrapf(jsonErr, "failed to read movie Provider configuration")
	}
	return movieProviders, nil
}

func getProviderMovies(movieProvider string, movieMap map[string][]*models.MovieItem) error {
	url := movieDataEndpoint + "/" + movieProvider + "/movies"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", movieDataApiKey)

	mRes, mErr := client.Do(req)
	if mErr != nil {
		return mErr
	}
	defer mRes.Body.Close()

	if mRes.StatusCode != 200 {
		// Read cached results from DynamoDB
		movieItems, rErr := repo.GetProviderMovies(movieTable, movieProvider)
		if rErr != nil {
			log.Err(rErr).Msgf("failed to read movies from respository for movie Provider %s", movieProvider)
			return rErr
		}
		movieMap[movieProvider] = models.SetReliable(false, movieItems)
		return nil
	}

	// Unmarshal the response, or if that fails try to read the most recent results for this
	// provider from DynamoDB. If that fails then we don't return any movies for that provider.
	var movieResponse models.MoviesResponse
	jErr := json.NewDecoder(mRes.Body).Decode(&movieResponse)
	if jErr != nil {
		log.Err(jErr).Msgf("failure to decode response from %s", url)

		movieItems, rErr := readCacheForProvider(movieProvider)
		if rErr != nil {
			return rErr
		}
		movieMap[movieProvider] = models.SetReliable(false, movieItems)
		return nil
	}
	movieMap[movieProvider] = models.SetReliable(true, movieResponse.Movies)
	return nil
}

func updateCacheForProvider(movieProvider string, movieItems []*models.MovieItem) error {
	log.Info().Msgf("updating cache for movie Provider %s", movieProvider)

	uErr := repo.UpdateProviderMovies(movieTable, movieProvider, movieItems)
	if uErr != nil {
		return uErr
	}
	return nil
}

func readCacheForProvider(movieProvider string) ([]*models.MovieItem, error) {
	log.Info().Msgf("reading cached results for movie Provider %s", movieProvider)

	movieItems, rErr := repo.GetProviderMovies(movieTable, movieProvider)
	if rErr != nil {
		return nil, rErr
	}
	return movieItems, nil
}

func main() {
	lambda.Start(handler)
}
