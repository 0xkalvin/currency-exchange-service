const pino = require('pino');
const { context } = require('./context');

const {
  APP_NAME,
  APP_TYPE,
  LOG_LEVEL,
  NODE_ENV,
} = process.env;

module.exports = (name) => pino({
  name,
  level: LOG_LEVEL,
  formatters: {
    level: (label) => ({ level: label }),
  },
  messageKey: 'message',
  timestamp: pino.stdTimeFunctions.isoTime,
  mixin() {
    const data = {
      app_name: APP_NAME,
      app_type: APP_TYPE,
      env: NODE_ENV,
    };

    const store = context.getStore();

    if (!store) {
      return data;
    }

    const customerId = store.get('customerId');
    const requestId = store.get('requestId');

    if (customerId) {
      data.customer_id = customerId;
    }

    if (requestId) {
      data.request_id = requestId;
    }

    return data;
  },
});
