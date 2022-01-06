const { BigNumber } = require('bignumber.js');

const movementRepository = require('../repositories/movement');
const orderRepository = require('../repositories/order');
const errors = require('../utils/errors');
const logger = require('../utils/logger')('MOVEMENT_SERVICE');

const exchangeIdentifier = 'CoolExchange';
const zero = new BigNumber(0);

async function createMovement(payloads) {
  const balanceUpdates = payloads.map((movement) => ({
    amount: movement.amount,
    currency_id: movement.currency_id,
    owner_id: movement.owner_id ? movement.owner_id : exchangeIdentifier,
    type: 'available',
  }));

  const movementsToCreate = payloads.map((movement) => ({
    id: movement.id,
    amount: movement.amount,
    currency_id: movement.currency_id,
    owner_id: movement.owner_id ? movement.owner_id : exchangeIdentifier,
    type: movement.type,
    source_id: movement.source_id,
  }));

  const sourceBalanceMovement = balanceUpdates[0];
  let movements;
  let orderTargetStatus = 'settled';

  try {
    const currentSourceCurrencyBalance = await movementRepository.getBalance(
      sourceBalanceMovement.owner_id,
      sourceBalanceMovement.currency_id,
      sourceBalanceMovement.type,
    );

    if (!currentSourceCurrencyBalance) {
      logger.warn({
        message: 'Balance does not exist',
        order_id: movementsToCreate[0].source_id,
        owner_id: sourceBalanceMovement.owner_id,
        currency_id: sourceBalanceMovement.currency_id,
        type: sourceBalanceMovement.type,
      });

      return;
    }

    const amountToSubtract = new BigNumber(sourceBalanceMovement.amount);
    const currentSourceCurrencyAmount = new BigNumber(currentSourceCurrencyBalance.amount);

    if (!currentSourceCurrencyAmount.plus(amountToSubtract).isGreaterThanOrEqualTo(zero)) {
      orderTargetStatus = 'failed';

      logger.debug({
        message: 'Insufficient balance to exchange',
        order_id: movementsToCreate[0].source_id,
        owner_id: sourceBalanceMovement.owner_id,
        currency_id: sourceBalanceMovement.currency_id,
        type: sourceBalanceMovement.type,
      });
    } else {
      movements = await movementRepository.createMovementsAndUpdateBalances(
        movementsToCreate,
        balanceUpdates,
      );

      logger.debug({
        message: 'Movements created successfully',
        order_id: movements[0].source_id,
      });
    }
  } catch (error) {
    if (error instanceof errors.InsufficientBalanceError) {
      orderTargetStatus = 'failed';
    } else {
      throw error;
    }
  }

  await orderRepository.enqueueOrderToSettle({
    customerId: movementsToCreate[0].owner_id,
    id: movementsToCreate[0].source_id,
    targetStatus: orderTargetStatus,
  });
}

module.exports = {
  createMovement,
};
