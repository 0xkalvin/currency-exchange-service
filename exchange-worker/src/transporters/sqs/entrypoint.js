const timers = require('timers/promises');

const {
  dynamodb,
  redis,
  sqs,
} = require('../../data-sources');
const logger = require('../../utils/logger')('WORKER_ENTRYPOINT');
const sqsConfig = require('../../config/sqs');
const SQSPoller = require('../../utils/sqs-poller');

const { createMovement } = require('./workers/movement-creation');
const { createOrder, rateLimitCreateOrder } = require('./workers/order-creation');
const { settleOrder } = require('./workers/order-settlement');

const queueToWorkerMap = new Map([
  ['movement-creation-queue', {
    eachMessage: createMovement,
  }],
  ['order-creation-queue', {
    beforePoll: rateLimitCreateOrder,
    eachMessage: createOrder,
  }],
  ['order-settlement-queue', {
    eachMessage: settleOrder,
  }],
]);

const pollers = [];

const {
  LOG_LEVEL,
} = process.env;

async function gracefullyShutdown() {
  try {
    await Promise.all(pollers.map((poller) => poller.stop()));

    logger.info({
      message: 'Cleanup has finished and process is about to shutdown',
      uptime: process.uptime(),
    });

    process.exit(0);
  } catch (error) {
    logger.error({
      message: 'Failed to gracefully shutdown worker',
      error_message: error.message,
      error_stack: error.stack,
    });
  }
}

async function run() {
  try {
    await Promise.all([
      dynamodb.checkConnection(),
      redis.checkConnection(),
      sqs.checkConnection(),
    ]);
  } catch (error) {
    logger.fatal({
      message: 'Failed to connect to data sources, exiting now',
      error_message: error.message,
      error_stack: error.stack,
    });

    process.exit(1);
  }

  logger.info({
    message: 'Starting exchange worker',
    worker_log_level: LOG_LEVEL,
  });

  sqsConfig.workerQueues.forEach((queueConfig) => {
    const worker = queueToWorkerMap.get(queueConfig.name);

    if (!worker) {
      logger.warn({
        message: 'There is no worker for queue',
        queue_name: queueConfig.name,
      });

      return;
    }

    const poller = new SQSPoller({
      queueUrl: queueConfig.queueURL,
      sqsClient: sqs.client,
    });

    logger.info({
      message: `Creating a worker for queue ${queueConfig.name}`,
    });

    poller.start({
      beforePoll: worker.beforePoll ? worker.beforePoll(poller, queueConfig) : null,
      eachMessage: worker.eachMessage,
    });

    poller.on('error', (error) => {
      logger.error({
        message: 'Failed to process message',
        queue_name: queueConfig.name,
        error_message: error.message,
        error_stack: error.stack,
      });
    });

    pollers.push(poller);
  });

  process.once('SIGTERM', async () => {
    await gracefullyShutdown();
  });
  process.once('SIGINT', async () => {
    await gracefullyShutdown();
  });
}

run();
