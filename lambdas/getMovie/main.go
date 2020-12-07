package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"movie-info/internal/utils"
)

var (
// Should be reading these from SSM
// movieDataEndpoint  = os.Getenv("MOVIE_DATA_ENDPOINT")
// movieDataApiKey    = os.Getenv("MOVIE_DATA_API_KEY")
// movieProviderNames = os.Getenv("MOVIE_PROVIDERS")
// movieTable         = os.Getenv("MOVIE_TABLE")
//
// client *http.Client
)

func init() {
	// client = &http.Client{}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("RequestID=%s, RequestTime=%s, Path=%s, PathParameters=%+v, QueryStringParameters=%+v",
		request.RequestContext.RequestID,
		request.RequestContext.RequestTime,
		request.Path,
		request.PathParameters,
		request.QueryStringParameters)
	/*
		var movieProviders []string
		jsonErr := json.Unmarshal([]byte(movieProviderNames), &movieProviders)
		if jsonErr != nil {
			return utils.CreateApiGwResponse(500,
				fmt.Sprintf("MOVIE_PROVIDERS environment variable is invalid. Value=%s", movieProviders)), nil
		}

		// movies := make(map[string][]*models.MovieItem)
		movieMap := make(map[string]*models.MovieItem)

		for _, movieProvider := range movieProviders {
			req, _ := http.NewRequest("GET", movieDataEndpoint+"/"+movieProvider+"/movies", nil)
			req.Header.Add("x-api-key", movieDataApiKey)
			mRes, mErr := client.Do(req)
			if mErr != nil {
				return utils.CreateApiGwResponse(500, mErr.Error()), nil
			}
			defer mRes.Body.Close()

			if mRes.StatusCode != 200 {
				// Read cached results from DynamoDB
				// movieItems, rErr := repository.GetProviderMovies(movieTable, movieProvider)
				// if rErr != nil {
				// 	log.Fatal(rErr)
				// 	continue
				// }
				// movies[movieProvider] = setReliable(false, movieItems)
				// aggregate(movies[movieProvider], movieMap)
				continue
			}

			// Unmarshal the response, or if that fails try to read the most recent results for this
			// provider from DynamoDB. If that fails then we don't return any movies for that provider.
			var movieResponse models.MoviesResponse
			jErr := json.NewDecoder(mRes.Body).Decode(&movieResponse)
			if jErr != nil {
				// Read cached results from DynamoDB
				// movieItems, rErr := repository.GetProviderMovies(movieTable, movieProvider)
				// if rErr != nil {
				// 	return utils.CreateApiGwResponse(500, jErr.Error()), nil
				// }
				// movies[movieProvider] = setReliable(false, movieItems)
				// aggregate(movies[movieProvider], movieMap)
				continue
			}
			uErr := repository.UpdateProviderMovies(movieTable, movieProvider, movieResponse.Movies)
			if uErr != nil {
				log.Fatalf(errors.Wrapf(uErr, "failed to update database").Error())
			}
			// movies[movieProvider] = setReliable(true, movieResponse.Movies)

			// aggregate(movies[movieProvider], movieMap)
		}

		payload, jsonErr := json.Marshal(movieMap)
		if jsonErr != nil {
			return utils.CreateApiGwResponse(500, jsonErr.Error()), nil
		}
		// OK, return the Movies
		return utils.CreateApiGwResponse(200, string(payload)), nil

	*/
	return utils.CreateApiGwResponse(200, "Hello"), nil
}

func main() {
	lambda.Start(handler)
}
