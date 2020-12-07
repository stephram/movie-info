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

func UpdateProviderMovies(tableName string, movieProvider string, movieItems []*models.MovieItem) error {
	ddb := createClient()

	for _, movieItem := range movieItems {
		dbMovieItem := convertToDbMovieItem(movieProvider, *movieItem)

		av, err := dynamodbattribute.MarshalMap(dbMovieItem)
		if err != nil {
			return errors.Wrap(err, "failed to marshal DbMovieItem")
		}
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		_, pErr := ddb.PutItem(input)
		if pErr != nil {
			return errors.Wrap(pErr, "failed to update database")
		}
	}
	log.Printf("UpdateProviderMovies: movieProvider=%s, count=%d", movieProvider, len(movieItems))
	return nil
}

func GetProviderMovies(tableName string, movieProvider string) ([]*models.MovieItem, error) {
	ddb := createClient()

	res, qErr := ddb.Query(&dynamodb.QueryInput{
		TableName: aws.String(tableName),
		// IndexName: aws.String("Provider"),
		KeyConditions: map[string]*dynamodb.Condition{
			"Provider": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(movieProvider),
					},
				},
			},
		},
	})
	if qErr != nil {
		return nil, errors.Wrapf(qErr, "failed to query for provider %s", movieProvider)
	}
	dbMovieItems := make([]DbMovieItem, len(res.Items))
	uErr := dynamodbattribute.UnmarshalListOfMaps(res.Items, &dbMovieItems)
	if uErr != nil {
		return nil, errors.Wrapf(uErr, "failed to unmarshal items")
	}
	log.Printf("GetProviderMovies: movieProvider=%s, count=%d", movieProvider, len(dbMovieItems))
	return convertToMovieItems(dbMovieItems), nil
}

func createClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	return svc
}
