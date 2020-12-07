package utils

import "github.com/aws/aws-lambda-go/events"

func CreateApiGwResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}
}
