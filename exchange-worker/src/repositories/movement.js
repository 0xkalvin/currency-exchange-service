const { TransactWriteItemsCommand } = require('@aws-sdk/client-dynamodb');
const { marshall } = require('@aws-sdk/util-dynamodb');

const { dynamodb, sqs } = require('../data-sources');
const logger = require('../utils/logger')('BALANCE_REPOSITORY');

const {
  BALANCE_TABLE,
  SQS_MOVEMENT_CREATION_QUEUE_URL,
} = process.env;

async function createMovementsAndUpdateBalances(movements, balanceUpdates) {
  try {
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
            pk: `BALANCE#${balanceOwnerId}#${balanceCurrencyId}#${balanceType}`,
            sk: 'BALANCE',
          }),
          UpdateExpression: 'SET amount = amount + :increment',
          ExpressionAttributeValues: {
            ':increment': {
              N: balanceIncrement,
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
        },
      });
    }

    await dynamodb.client.send(new TransactWriteItemsCommand({
      TransactItems: transactItems,
    }));

    return movements;
  } catch (error) {
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

module.exports = {
  createMovementsAndUpdateBalances,
  enqueueMovements,
};
