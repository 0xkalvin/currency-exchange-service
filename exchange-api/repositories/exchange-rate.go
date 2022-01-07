package repositories

import (
	"context"
	"exchange-api/entities"
	"exchange-api/utils/logger"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (
	ExchangeRateRepositoryInterface interface {
		CreateExchangeRate(ctx context.Context, payload *entities.ExchangeRate) (*entities.ExchangeRate, error)
	}

	ExchangeRateRepository struct {
		dynamoClient *dynamodb.DynamoDB
	}
)

var exchangeRateLogger = logger.NewLogger()
var exchangeTableName = os.Getenv("EXCHANGE_TABLE")

func NewExchangeRateRepository(dynamo *dynamodb.DynamoDB) ExchangeRateRepositoryInterface {
	return ExchangeRateRepository{
		dynamoClient: dynamo,
	}
}

func (r ExchangeRateRepository) CreateExchangeRate(ctx context.Context, exchangeRate *entities.ExchangeRate) (*entities.ExchangeRate, error) {
	requestId := ctx.Value("requestId").(string)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(exchangeRate.Id),
			},
			"sk": {
				S: aws.String("EXCHANGE_RATE"),
			},
			"currency_id": {
				S: aws.String(exchangeRate.CurrencyId),
			},
			"base_currency_id": {
				S: aws.String(exchangeRate.BaseCurrencyId),
			},
			"rate": {
				S: aws.String(exchangeRate.Rate),
			},
			"timestamp": {
				S: aws.String(exchangeRate.Timestamp),
			},
		},
		TableName: aws.String(exchangeTableName),
	}

	_, err := r.dynamoClient.PutItem(input)

	if err != nil {
		exchangeRateLogger.Error(map[string]interface{}{
			"message":       "Failed to create exchange in dynamodb",
			"error_message": err.Error(),
			"request_id":    requestId,
		})

		return nil, err
	}

	return exchangeRate, nil
}
