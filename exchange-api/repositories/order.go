package repositories

import (
	"context"
	"encoding/json"
	"exchange-api/entities"
	"exchange-api/utils/logger"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type (
	OrderRepositoryInterface interface {
		CreateOrder(ctx context.Context, orderPayload *entities.Order) (*entities.Order, error)
		GetOrdersDashboardByCustomer(ctx context.Context, customerId string) ([]*OrderStatusTotalItem, error)
	}

	OrderRepository struct {
		sqsClient    *sqs.SQS
		dynamoClient *dynamodb.DynamoDB
	}

	OrderCustomerItem struct {
		Pk    string
		Sk    string
		Total string
	}

	OrderStatusTotalItem struct {
		Status string
		Total  string
	}
)

var orderLogger = logger.NewLogger()
var sqsOrderCreationQueueURL = os.Getenv("SQS_ORDER_CREATION_QUEUE_URL")

func NewOrderRepository(sqs *sqs.SQS, dynamoClient *dynamodb.DynamoDB) OrderRepositoryInterface {
	return OrderRepository{
		sqsClient:    sqs,
		dynamoClient: dynamoClient,
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
		QueueUrl:               aws.String(sqsOrderCreationQueueURL),
		MessageGroupId:         aws.String(orderPayload.CustomerId),
		MessageDeduplicationId: aws.String(orderPayload.Id),
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	_, err = r.sqsClient.SendMessage(input)

	if err != nil {
		orderLogger.Error(map[string]interface{}{
			"message":       "Failed to enqueue order into SQS",
			"error_message": err.Error(),
			"request_id":    requestId,
		})

		return nil, err
	}

	return orderPayload, nil
}

func (r OrderRepository) GetOrdersDashboardByCustomer(ctx context.Context, customerId string) ([]*OrderStatusTotalItem, error) {
	requestId := ctx.Value("requestId").(string)

	output, err := r.dynamoClient.Query(&dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("EXCHANGE_TABLE")),
		KeyConditions: map[string]*dynamodb.Condition{
			"pk": {
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(customerId),
					},
				},
				ComparisonOperator: aws.String("EQ"),
			},
			"sk": {
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(customerId),
					},
				},
				ComparisonOperator: aws.String("BEGINS_WITH"),
			},
		},
	})

	if err != nil {
		orderLogger.Error(map[string]interface{}{
			"message":       "Failed to get order_customer entity in dynamodb",
			"error_message": err.Error(),
			"request_id":    requestId,
		})

		return nil, err
	}

	items := []OrderCustomerItem{}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	if err != nil {
		return nil, err
	}

	responseItems := []*OrderStatusTotalItem{}

	for _, item := range items {
		status := strings.Split(item.Sk, "#")[2]

		responseItems = append(responseItems, &OrderStatusTotalItem{
			Status: status,
			Total:  item.Total,
		})
	}

	return responseItems, nil
}
