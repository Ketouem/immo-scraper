package db

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sirupsen/logrus"

	"github.com/Ketouem/immo-scraper/lib/scraper"
)

var db *dynamodb.DynamoDB

// Setup prepares the connection handler
func Setup(endpointURL string) {
	config := aws.NewConfig()
	if len(endpointURL) != 0 {
		config = config.WithEndpoint(endpointURL)
	}
	session := session.Must(session.NewSession())
	db = dynamodb.New(session, config)
}

// Get returns the connection handler
func Get() (*dynamodb.DynamoDB, error) {
	if db == nil {
		return nil, errors.New("Setup needs to be called first")
	}
	return db, nil
}

// Provision creates the needed tables and attributes
func Provision(db *dynamodb.DynamoDB) (err error) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("link"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("link"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Results"),
	}
	_, err = db.CreateTable(input)
	if err != nil {
		if serr, ok := err.(awserr.Error); ok {
			if serr.Code() == dynamodb.ErrCodeResourceInUseException {
				err = nil
			}
		}
	}
	return
}

// PutResults store results into db
func PutResults(database *dynamodb.DynamoDB, results []scraper.Result) (err error) {
	for _, result := range results {
		rs, _ := dynamodbattribute.MarshalMap(result)
		input := &dynamodb.PutItemInput{
			Item:                rs,
			TableName:           aws.String("Results"),
			ConditionExpression: aws.String("attribute_not_exists(link)"),
		}
		_, err := database.PutItem(input)
		if err != nil {
			if serr, ok := err.(awserr.Error); ok {
				// Raised whent the object already exists in db
				if serr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
					logrus.WithField("link", result.Link).Debug("Result already exists, skipping peristence.")
					continue
				} else {
					return err
				}
			}
		} else {
			logrus.WithField("link", result.Link).Debug("New result, persisting.")
		}
	}
	return
}

// FetchNewResults retrieves unnotified results from the database
func FetchNewResults(database *dynamodb.DynamoDB) (results []scraper.Result, err error) {
	notNotifiedFilter := expression.Name("Notified").Equal(expression.Value(false))
	expr, err := expression.NewBuilder().WithFilter(notNotifiedFilter).Build()
	params := &dynamodb.ScanInput{
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String("Results"),
	}
	result, err := database.Scan(params)
	for _, i := range result.Items {
		rs := scraper.Result{}
		err = dynamodbattribute.UnmarshalMap(i, &rs)
		results = append(results, rs)
	}
	return
}
