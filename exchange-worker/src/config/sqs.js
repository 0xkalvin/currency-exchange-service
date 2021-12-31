module.exports = {
  movementCreationQueueURL: process.env.SQS_MOVEMENT_CREATION_QUEUE_URL,
  orderCreationQueueURL: process.env.SQS_ORDER_CREATION_QUEUE_URL,
  orderSettlementQueueURL: process.env.SQS_ORDER_SETTLEMENT_QUEUE_URL,
  endpoint: process.env.SQS_ENDPOINT,
  maxRetries: process.env.SQS_MAX_RETRIES,
  region: process.env.SQS_REGION,
  workerQueues: [
    {
      concurrency: process.env.SQS_CONCURRENCY,
      queueURL: process.env.SQS_MOVEMENT_CREATION_QUEUE_URL,
      name: 'movement-creation-queue',
    },
    {
      concurrency: process.env.SQS_CONCURRENCY,
      queueURL: process.env.SQS_ORDER_CREATION_QUEUE_URL,
      name: 'order-creation-queue',
    },
    {
      concurrency: process.env.SQS_CONCURRENCY,
      queueURL: process.env.SQS_ORDER_SETTLEMENT_QUEUE_URL,
      name: 'order-settlement-queue',
    },
  ],
};
