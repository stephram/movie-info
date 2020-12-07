package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"log"
	"movie-info/internal/models"
)

func UpdateProviderMovies(tableName string, movieProvider string, movieItems []models.MovieItem) error {
	ddb := createClient()

	for _, movieItem := range movieItems {
		av, err := dynamodbattribute.MarshalMap(movieItem)
		if err != nil {
			return errors.Wrap(err, "failed to marshal MovieItem")
		}
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		_, pErr := ddb.PutItem(input)
		if pErr != nil {
			return errors.Wrap(pErr, "failed to update database")
		}
		log.Printf("Added Movie, ID=%s, Title=%s", movieItem.ID, movieItem.Title)
	}
	return nil
}

func ReadProviderMovies(tableName string, provider string) ([]models.MovieItem, error) {
	ddb := createClient()

	result, gErr := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"provider": {
				N: aws.String(provider),
			},
		},
	})
	if gErr != nil {
		return nil, errors.Wrap(gErr, "failed to GetItem")
	}

	movieItems := make([]models.MovieItem, 0)

	uErr := dynamodbattribute.UnmarshalMap(result.Item, &movieItems)
	if uErr != nil {
		return nil, errors.Wrap(uErr, "failed to unmarshal movies for provider "+provider)
	}
	return movieItems, nil
}

func createClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	return svc
}
