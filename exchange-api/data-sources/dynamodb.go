package datasources

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func BuildDynamoDBClient() *dynamodb.DynamoDB {
	session := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Endpoint: aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
				Region:   aws.String(os.Getenv("DYNAMODB_REGION")),
			},
		},
	))

	dynamoDBClient := dynamodb.New(session)

	return dynamoDBClient
}
