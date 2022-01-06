const crypto = require('crypto');
const { BigNumber } = require('bignumber.js');

const movementRepository = require('../repositories/movement');
const exchangeRateRepository = require('../repositories/exchange-rate');
const orderRepository = require('../repositories/order');
const logger = require('../utils/logger')('ORDER_SERVICE');

const feeRate = new BigNumber('0.005');

const orderStateMachine = new Map([
  ['created', new Set(['settled', 'failed'])],
  ['settled', new Set([])],
  ['failed', new Set([])],
]);

async function createOrder(payload) {
  const {
    amount,
    customerId,
    id,
    sourceCurrencyId,
    scheduleDate,
    targetCurrencyId,
  } = payload;

  const exchangeRatesResult = await Promise.all([
    exchangeRateRepository.findExchangeRateById(sourceCurrencyId),
    exchangeRateRepository.findExchangeRateById(targetCurrencyId),
  ]);

  if (!exchangeRatesResult[0] || !exchangeRatesResult[1]) {
    logger.warn({
      message: 'Exchange rate does not exist for specified currency',
      currency_id: exchangeRatesResult[0] ? targetCurrencyId : sourceCurrencyId,
    });

    return null;
  }

  const order = await orderRepository.createOrder({
    amount,
    customerId,
    id,
    sourceCurrencyId,
    scheduleDate,
    targetCurrencyId,
    status: 'created',
  });

  if (scheduleDate) {
    logger.info({
      message: 'Order has been schedule for future processing',
      schedule_date: scheduleDate,
    });

    return order;
  }

  const sourceAmount = new BigNumber(amount);

  const { base_currency_id: baseCurrencyId } = exchangeRatesResult[0];

  const sourceCurrencyRate = new BigNumber(exchangeRatesResult[0].rate);
  const targetCurrencyRate = new BigNumber(exchangeRatesResult[1].rate);

  const convertedAmount = sourceAmount
    .multipliedBy(sourceCurrencyRate)
    .dividedBy(targetCurrencyRate);

  const feeAmount = convertedAmount.multipliedBy(feeRate);
  const targetAmount = convertedAmount.minus(feeAmount);

  const movements = [
    {
      id: crypto.randomUUID(),
      amount: sourceAmount.negated().toString(),
      type: 'exchange',
      currency_id: sourceCurrencyId,
      owner_id: customerId,
      source_id: id,
    },
    {
      id: crypto.randomUUID(),
      amount: targetAmount.toString(),
      type: 'exchange',
      currency_id: targetCurrencyId,
      owner_id: customerId,
      source_id: id,
    },
    {
      id: crypto.randomUUID(),
      amount: feeAmount.toString(),
      type: 'fee',
      currency_id: baseCurrencyId,
      source_id: id,
    },
  ];

  await movementRepository.enqueueMovements(movements);

  logger.debug({
    message: 'Order created successfully',
    order_id: id,
  });

  return order;
}

async function settleOrder(payload) {
  const { customerId, id, targetStatus } = payload;

  const order = await orderRepository.findOrderById(id);

  if (!order) {
    logger.warn({
      message: 'Order does not exist',
      order_id: id,
    });

    return null;
  }

  const possibleTransitions = orderStateMachine.get(order.status);

  if (!possibleTransitions || !possibleTransitions.has(targetStatus)) {
    logger.warn({
      message: 'Order cannot transition to specified status',
      current_status: order.status,
      target_status: targetStatus,
    });

    return null;
  }

  const updatedOrder = await orderRepository.settleOrder({
    customerId,
    id,
    targetStatus,
  });

  logger.debug({
    message: 'Order settled successfully',
    new_status: targetStatus,
    order_id: id,
  });

  return updatedOrder;
}

module.exports = {
  createOrder,
  settleOrder,
};
