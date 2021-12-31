const orderService = require('../../../services/order');

function createOrder(queueConfig) {
  return async function innerCreateOrder(message) {
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
  };
}

module.exports = {
  createOrder,
};
