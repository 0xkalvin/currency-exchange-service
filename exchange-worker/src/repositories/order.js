const {
  GetItemCommand, TransactWriteItemsCommand,
} = require('@aws-sdk/client-dynamodb');
const { marshall, unmarshall } = require('@aws-sdk/util-dynamodb');

const { dynamodb, sqs } = require('../data-sources');
const sqsConfig = require('../config/sqs');
const logger = require('../utils/logger')('ORDER_REPOSITORY');

const {
  EXCHANGE_TABLE,
} = process.env;

async function createOrder(payload) {
  const {
    amount,
    customerId,
    id,
    status,
    sourceCurrencyId,
    scheduleDate,
    targetCurrencyId,
  } = payload;

  try {
    const item = {
      pk: id,
      sk: 'ORDER',
      amount,
      customer_id: customerId,
      status,
      source_currency_id: sourceCurrencyId,
      target_currency_id: targetCurrencyId,
    };

    if (scheduleDate) {
      item.schedule_date = scheduleDate;
    }

    await dynamodb.client.send(new TransactWriteItemsCommand({
      TransactItems: [
        {
          Put: {
            TableName: EXCHANGE_TABLE,
            Item: marshall(item),
            ConditionExpression: 'attribute_not_exists(#pk) AND attribute_not_exists(#sk)',
            ExpressionAttributeNames: {
              '#pk': 'pk',
              '#sk': 'sk',
            },
          },
        },
        {
          Update: {
            TableName: EXCHANGE_TABLE,
            Key: marshall({
              pk: `${customerId}`,
              sk: `${customerId}#ORDER_CUSTOMER#${status}`,
            }),
            UpdateExpression: 'SET #total = #total + :inc',
            ExpressionAttributeValues: {
              ':inc': {
                N: '1',
              },
            },
            ExpressionAttributeNames: {
              '#total': 'total',
            },
          },
        },
      ],
    }));
  } catch (error) {
    logger.error({
      message: 'Failed to create order on dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

async function enqueueOrderToSettle(payload) {
  const {
    customerId,
    id,
    targetStatus,
  } = payload;

  try {
    await sqs.client.sendMessage({
      QueueUrl: sqsConfig.orderSettlementQueueURL,
      MessageBody: JSON.stringify({
        customer_id: customerId,
        id,
        target_status: targetStatus,
      }),
      MessageGroupId: customerId,
      MessageDeduplicationId: id,
    });
  } catch (error) {
    logger.error({
      message: 'Failed to enqueue order to settle onto sqs',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

async function findOrderById(id) {
  try {
    const result = await dynamodb.client.send(new GetItemCommand({
      TableName: EXCHANGE_TABLE,
      Key: marshall({
        pk: id,
        sk: 'ORDER',
      }),
    }));

    return result.Item ? unmarshall(result.Item) : null;
  } catch (error) {
    logger.error({
      message: 'Failed to find order on dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

async function settleOrder(payload) {
  const {
    id,
    customerId,
    targetStatus,
  } = payload;

  try {
    await dynamodb.client.send(new TransactWriteItemsCommand({
      TransactItems: [
        {
          Update: {
            TableName: EXCHANGE_TABLE,
            Key: marshall({
              pk: id,
              sk: 'ORDER',
            }),
            UpdateExpression: 'SET #currentStatus = :targetStatus',
            ExpressionAttributeValues: marshall({
              ':targetStatus': targetStatus,
            }),
            ExpressionAttributeNames: {
              '#currentStatus': 'status',
            },
          },
        },
        {
          Update: {
            TableName: EXCHANGE_TABLE,
            Key: marshall({
              pk: `${customerId}`,
              sk: `${customerId}#ORDER_CUSTOMER#${targetStatus}`,
            }),
            UpdateExpression: 'SET #total = #total + :inc',
            ExpressionAttributeValues: {
              ':inc': {
                N: '1',
              },
            },
            ExpressionAttributeNames: {
              '#total': 'total',
            },
          },
        },
      ],
    }));
  } catch (error) {
    logger.error({
      message: 'Failed to update order on dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

module.exports = {
  createOrder,
  enqueueOrderToSettle,
  findOrderById,
  settleOrder,
};
