package repositories

import (
	"context"
	"exchange-api/utils/logger"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type (
	BalanceRepositoryInterface interface {
		GetMovementsDashboardByCustomer(ctx context.Context, customerId string) ([]*MovementTypeTotalItem, error)
		GetBalanceByCustomer(ctx context.Context, customerId string) ([]*BalanceResponseItem, error)
	}

	BalanceRepository struct {
		dynamoClient *dynamodb.DynamoDB
	}

	MovementCustomerItem struct {
		Pk    string
		Sk    string
		Total string
	}

	MovementTypeTotalItem struct {
		Type  string
		Total string
	}

	BalanceItem struct {
		Pk         string
		Sk         string
		Amount     string
		OwnerId    string
		CurrencyId string
	}

	BalanceResponseItem struct {
		Amount     string
		CurrencyId string
	}
)

var balanceLogger = logger.NewLogger()

func NewBalanceRepository(dynamoClient *dynamodb.DynamoDB) BalanceRepositoryInterface {
	return BalanceRepository{
		dynamoClient: dynamoClient,
	}
}

func (r BalanceRepository) GetMovementsDashboardByCustomer(ctx context.Context, customerId string) ([]*MovementTypeTotalItem, error) {
	requestId := ctx.Value("requestId").(string)

	output, err := r.dynamoClient.Query(&dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("BALANCE_TABLE")),
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
		balanceLogger.Error(map[string]interface{}{
			"message":       "Failed to get movement_customer entity in dynamodb",
			"error_message": err.Error(),
			"request_id":    requestId,
		})

		return nil, err
	}

	items := []MovementCustomerItem{}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	if err != nil {
		return nil, err
	}

	responseItems := []*MovementTypeTotalItem{}

	for _, item := range items {
		movementType := strings.Split(item.Sk, "#")[2]

		responseItems = append(responseItems, &MovementTypeTotalItem{
			Type:  movementType,
			Total: item.Total,
		})
	}

	return responseItems, nil
}

func (r BalanceRepository) GetBalanceByCustomer(ctx context.Context, customerId string) ([]*BalanceResponseItem, error) {
	requestId := ctx.Value("requestId").(string)

	output, err := r.dynamoClient.Query(&dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("BALANCE_TABLE")),
		IndexName: aws.String("owner_id_index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"owner_id": {
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
						S: aws.String("BALANCE"),
					},
				},
				ComparisonOperator: aws.String("EQ"),
			},
		},
	})

	if err != nil {
		balanceLogger.Error(map[string]interface{}{
			"message":       "Failed to get balance in dynamodb",
			"error_message": err.Error(),
			"request_id":    requestId,
		})

		return nil, err
	}

	items := []BalanceItem{}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	if err != nil {
		return nil, err
	}

	responseItems := []*BalanceResponseItem{}

	for _, item := range items {
		responseItems = append(responseItems, &BalanceResponseItem{
			Amount:     item.Amount,
			CurrencyId: item.CurrencyId,
		})
	}

	return responseItems, nil
}
