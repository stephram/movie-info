package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("RequestID=%s, Path=%s, IP=%v, MOVIE_DATA_ENDPOINT=%s, MOVIE_DATA_API_KEY=%s, MOVIE_PROVIDERS=%v, MOVIE_TABLE=%s",
			request.RequestContext.RequestID,
			request.Path,
			string(ip),
			os.Getenv("MOVIE_DATA_ENDPOINT"),
			os.Getenv("MOVIE_DATA_API_KEY"),
			os.Getenv("MOVIE_PROVIDERS"),
			os.Getenv("MOVIE_TABLE"),
		),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
