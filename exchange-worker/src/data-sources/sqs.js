const {
  SQS,
  ListQueueTagsCommand,
} = require('@aws-sdk/client-sqs');

const logger = require('../utils/logger')('SQS_INDEX');
const config = require('../config/sqs');

const client = new SQS({
  endpoint: config.endpoint,
  maxRetries: config.maxRetries,
  region: config.region,
});

async function checkConnection() {
  try {
    await Promise.all([
      client.send(new ListQueueTagsCommand({
        QueueUrl: config.movementCreationQueueURL,
      })),
      client.send(new ListQueueTagsCommand({
        QueueUrl: config.orderCreationQueueURL,
      })),
      client.send(new ListQueueTagsCommand({
        QueueUrl: config.orderSettlementQueueURL,
      })),
    ]);

    logger.info({
      message: 'Successfully connected to sqs queues',
    });
  } catch (error) {
    logger.error({
      message: 'Failed to connect to sqs',
      error_message: error.message,
      error_stack: error.stack,
    });

    throw error;
  }
}

module.exports = {
  client,
  checkConnection,
};
