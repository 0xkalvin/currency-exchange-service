const { DynamoDBClient, DescribeTableCommand } = require('@aws-sdk/client-dynamodb');

const config = require('../../config/dynamodb');
const logger = require('../../utils/logger')('DYNAMODB_INDEX');

const client = new DynamoDBClient(config);

const {
  BALANCE_TABLE,
  EXCHANGE_TABLE,
} = process.env;

async function checkConnection() {
  try {
    await Promise.all([
      client.send(new DescribeTableCommand({
        TableName: BALANCE_TABLE,
      })),
      client.send(new DescribeTableCommand({
        TableName: EXCHANGE_TABLE,
      })),
    ]);

    logger.info({
      message: 'Successfully connected to dynamoDB',
    });
  } catch (error) {
    logger.error({
      message: 'Failed to connect to dynamodb',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

module.exports = {
  checkConnection,
  client,
};
