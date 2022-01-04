const { RateLimiterRedis, RateLimiterRes } = require('rate-limiter-flexible');

const { redis } = require('../../../data-sources');
const orderService = require('../../../services/order');
const logger = require('../../../utils/logger')('ORDER_CREATION_WORKER');

function rateLimitCreateOrder(poller, queueConfig) {
  const rateLimiter = new RateLimiterRedis({
    storeClient: redis.client,
    points: queueConfig.rateLimit.points,
    duration: queueConfig.rateLimit.duration,
  });

  return async function innerRateLimitCreateOrder() {
    try {
      const { lastMessagesCount } = poller;

      if (lastMessagesCount > 0) {
        await rateLimiter.consume(queueConfig.name, 1);
      }
    } catch (error) {
      if (error instanceof RateLimiterRes) {
        const { msBeforeNext } = error;

        logger.warn({
          message: `Queue ${queueConfig.name} hit rate limit`,
          points: queueConfig.rateLimit.points,
          duration: queueConfig.rateLimit.duration,
          ms_before_next: msBeforeNext,
        });

        if (poller.isRunning) {
          await poller.stop();

          setTimeout(() => {
            poller.resume();
          }, msBeforeNext);
        }
      } else {
        logger.error({
          message: 'Failed to apply rate limit. Continuing consumption at full speed',
          queue_name: queueConfig.name,
          error_message: error.message,
          error_stack: error.stack,
        });
      }
    }
  };
}

async function createOrder(message) {
  const payload = JSON.parse(message.Body);

  const {
    amount,
    customer_id: customerId,
    id,
    source_currency_id: sourceCurrencyId,
    target_currency_id: targetCurrencyId,
  } = payload;

  await orderService.createOrder({
    amount,
    customerId,
    id,
    sourceCurrencyId,
    targetCurrencyId,
  });
}

module.exports = {
  createOrder,
  rateLimitCreateOrder,
};
