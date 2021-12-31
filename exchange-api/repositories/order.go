package repositories

import (
	"encoding/json"
	"exchange-api/entities"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type (
	OrderRepositoryInterface interface {
		CreateOrder(orderPayload *entities.Order) (*entities.Order, error)
	}

	OrderRepository struct {
		sqsClient *sqs.SQS
	}
)

var sqsOrderCreationQueueURL = os.Getenv("SQS_ORDER_CREATION_QUEUE_URL")

func NewOrderRepository(sqs *sqs.SQS) OrderRepositoryInterface {
	return OrderRepository{
		sqsClient: sqs,
	}
}

func (r OrderRepository) CreateOrder(orderPayload *entities.Order) (*entities.Order, error) {
	messageBody, err := json.Marshal(orderPayload)

	if err != nil {
		return nil, err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
		QueueUrl:    aws.String(sqsOrderCreationQueueURL),
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	_, err = r.sqsClient.SendMessage(input)

	if err != nil {
		return nil, err
	}

	return orderPayload, nil
}
