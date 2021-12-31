const { GetItemCommand } = require('@aws-sdk/client-dynamodb');
const { marshall, unmarshall } = require('@aws-sdk/util-dynamodb');

const { dynamodb } = require('../data-sources');
const logger = require('../utils/logger')('EXCHANGE_RATE_REPOSITORY');

const {
  EXCHANGE_TABLE,
} = process.env;

async function findExchangeRateById(id) {
  try {
    const result = await dynamodb.client.send(new GetItemCommand({
      TableName: EXCHANGE_TABLE,
      Key: marshall({
        pk: id,
        sk: 'EXCHANGE_RATE',
      }),
    }));

    return result.Item ? unmarshall(result.Item) : null;
  } catch (error) {
    logger.error({
      message: 'Failed to find exchange rate on dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

module.exports = {
  findExchangeRateById,
};
