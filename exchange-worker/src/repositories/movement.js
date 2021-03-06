const { GetItemCommand, TransactWriteItemsCommand } = require('@aws-sdk/client-dynamodb');
const { marshall, unmarshall } = require('@aws-sdk/util-dynamodb');

const { dynamodb, sqs } = require('../data-sources');
const errors = require('../utils/errors');
const logger = require('../utils/logger')('BALANCE_REPOSITORY');

const {
  BALANCE_TABLE,
  SQS_MOVEMENT_CREATION_QUEUE_URL,
} = process.env;

async function createMovementsAndUpdateBalances(movements, balanceUpdates) {
  const transactItems = [];

  for (let i = 0; i < movements.length; i += 1) {
    const {
      amount: balanceIncrement,
      currency_id: balanceCurrencyId,
      owner_id: balanceOwnerId,
      type: balanceType,
    } = balanceUpdates[i];

    transactItems.push({
      Update: {
        TableName: BALANCE_TABLE,
        Key: marshall({
          pk: `${balanceOwnerId}#${balanceCurrencyId}#${balanceType}`,
          sk: 'BALANCE',
        }),
        UpdateExpression: 'SET amount = amount + :increment',
        ConditionExpression: 'amount > :zero',
        ExpressionAttributeValues: {
          ':increment': {
            N: balanceIncrement,
          },
          ':zero': {
            N: '0',
          },
        },
      },
    });

    const {
      id: movementId,
      amount: movementAmount,
      currency_id: movementCurrencyId,
      owner_id: movementOwnerId,
      source_id: movementSourceId,
      type: movementType,
    } = movements[i];

    transactItems.push({
      Put: {
        TableName: BALANCE_TABLE,
        Item: marshall({
          pk: movementId,
          sk: 'MOVEMENT',
          amount: movementAmount,
          currency_id: movementCurrencyId,
          owner_id: movementOwnerId,
          source_id: movementSourceId,
          type: movementType,
        }),
        ConditionExpression: 'attribute_not_exists(#pk) AND attribute_not_exists(#sk)',
        ExpressionAttributeNames: {
          '#pk': 'pk',
          '#sk': 'sk',
        },
      },
    });
  }

  try {
    const ownerId = movements[0].owner_id;
    const movementType = movements[0].type;

    await dynamodb.client.send(new TransactWriteItemsCommand({
      TransactItems: [
        ...transactItems,
        {
          Update: {
            TableName: BALANCE_TABLE,
            Key: marshall({
              pk: `${ownerId}`,
              sk: `${ownerId}#MOVEMENT_OWNER#${movementType}`,
            }),
            UpdateExpression: 'SET #total = #total + :inc',
            ExpressionAttributeValues: {
              ':inc': {
                N: '2',
              },
            },
            ExpressionAttributeNames: {
              '#total': 'total',
            },
          },
        },
      ],
    }));

    return movements;
  } catch (error) {
    if (error.name === 'TransactionCanceledException' && error.$fault === 'client') {
      throw new errors.InsufficientBalanceError(
        `Owner ${balanceUpdates[0].owner_id} has insufficient funds for swapping currency ${balanceUpdates[0].currency_id}`,
      );
    }

    logger.error({
      message: 'Failed to create movements on dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

async function enqueueMovements(movements) {
  try {
    await sqs.client.sendMessage({
      QueueUrl: SQS_MOVEMENT_CREATION_QUEUE_URL,
      MessageBody: JSON.stringify(movements),
      MessageGroupId: movements[0].owner_id,
      MessageDeduplicationId: movements[0].source_id,
    });
  } catch (error) {
    logger.error({
      message: 'Failed to enqueue movements onto sqs',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

async function getBalance(ownerId, currencyId, type) {
  try {
    const pk = `${ownerId}#${currencyId}#${type}`;

    const result = await dynamodb.client.send(new GetItemCommand({
      TableName: BALANCE_TABLE,
      Key: marshall({
        pk,
        sk: 'BALANCE',
      }),
    }));

    return result.Item ? unmarshall(result.Item) : null;
  } catch (error) {
    logger.error({
      message: 'Failed to find balance on dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

module.exports = {
  createMovementsAndUpdateBalances,
  enqueueMovements,
  getBalance,
};
