package datasources

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func BuildSQSClient() *sqs.SQS {

	session := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Endpoint: aws.String(os.Getenv("SQS_ENDPOINT")),
				Region:   aws.String(os.Getenv("SQS_REGION")),
			},
		},
	))

	sqsClient := sqs.New(session)

	return sqsClient
}
