const Redis = require('ioredis');

const config = require('../config/redis');
const logger = require('../utils/logger')('REDIS_INDEX');

const client = new Redis(config);

let hasInitialized = false;

async function checkConnection() {
  if (!hasInitialized) {
    await new Promise((resolve) => {
      client.once('ready', function done() {
        client.removeListener('ready', done);

        resolve();
      });
    });

    hasInitialized = true;
  }

  await client.ping();

  logger.info({
    message: 'Successfully connected to redis',
  });
}

module.exports = {
  client,
  checkConnection,
};
