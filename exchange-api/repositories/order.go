package repositories

import (
	"context"
	"encoding/json"
	"exchange-api/entities"
	"exchange-api/utils/logger"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type (
	OrderRepositoryInterface interface {
		CreateOrder(ctx context.Context, orderPayload *entities.Order) (*entities.Order, error)
	}

	OrderRepository struct {
		sqsClient *sqs.SQS
	}
)

var log = logger.NewLogger()
var sqsOrderCreationQueueURL = os.Getenv("SQS_ORDER_CREATION_QUEUE_URL")

func NewOrderRepository(sqs *sqs.SQS) OrderRepositoryInterface {
	return OrderRepository{
		sqsClient: sqs,
	}
}

func (r OrderRepository) CreateOrder(ctx context.Context, orderPayload *entities.Order) (*entities.Order, error) {
	messageBody, err := json.Marshal(orderPayload)

	if err != nil {
		return nil, err
	}

	requestId := ctx.Value("requestId").(string)

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"x-request-id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(requestId),
			},
		},
		QueueUrl: aws.String(sqsOrderCreationQueueURL),
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	_, err = r.sqsClient.SendMessage(input)

	if err != nil {
		log.Error(map[string]interface{}{
			"message":       "Failed to enqueue order into SQS",
			"error_message": err.Error(),
			"request_id":    requestId,
		})

		return nil, err
	}

	return orderPayload, nil
}
