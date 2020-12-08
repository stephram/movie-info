package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"
	_ "log"
	"movie-info/internal/models"
	"movie-info/internal/utils"
)

var log = utils.GetLogger()

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
	log.Info().Msgf("movieProvider=%s, count=%d", movieProvider, len(movieItems))

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
		return nil, errors.Wrapf(uErr, "failed to unmarshal movieItems")
	}
	log.Info().Msgf("movieProvider=%s, count=%d", movieProvider, len(dbMovieItems))
	return convertToMovieItems(dbMovieItems), nil
}

func UpdateMovieItem(tableName string, movieItem *models.MovieItem) error {
	ddb := createClient()

	dbMovieItem := convertToDbMovieItem(movieItem.Provider, *movieItem)

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
		return errors.Wrapf(pErr, "failed to update movieItem ID %s", movieItem.ID)
	}
	log.Info().Interface("movieItem", movieItem).Msg("Database update successful")
	return nil
}

func GetMoviesByMovieID(tableName string, movieID string) ([]*models.MovieItem, error) {
	ddb := createClient()

	// Create the Expression to fill the input struct with.
	filt := expression.Name("MovieID").Equal(expression.Value(movieID))
	proj := expression.NamesList(
		expression.Name("MovieID"),
		expression.Name("ID"),
		expression.Name("Title"),
		expression.Name("Provider"),
		expression.Name("Price"),
		expression.Name("Type"),
		expression.Name("Poster"),
	)
	expr, eErr := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if eErr != nil {
		return nil, errors.Wrapf(eErr, "failed to create dynamodb Expression for movieID %s", movieID)
	}

	dRes, dErr := ddb.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	})
	if dErr != nil {
		return nil, errors.Wrapf(dErr, "failed to Scan for movieID %s", movieID)
	}
	dbMovieItems := make([]DbMovieItem, 0)

	uErr := dynamodbattribute.UnmarshalListOfMaps(dRes.Items, &dbMovieItems)
	if uErr != nil {
		return nil, errors.Wrapf(uErr, "failed to unmarshal DbMovieItems")
	}
	log.Info().Interface("dbMovieItems", dbMovieItems).Msgf("read dbMovieItems from %s", tableName)

	return convertToMovieItems(dbMovieItems), nil
}

func createClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	return dynamodb.New(sess)
}
