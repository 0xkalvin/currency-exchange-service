const movementRepository = require('../repositories/movement');
const orderRepository = require('../repositories/order');
const errors = require('../utils/errors');
const logger = require('../utils/logger')('MOVEMENT_SERVICE');

const exchangeIdentifier = 'CoolExchange';

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

  let movements;
  let orderTargetStatus = 'settled';

  try {
    movements = await movementRepository.createMovementsAndUpdateBalances(
      movementsToCreate,
      balanceUpdates,
    );

    logger.debug({
      message: 'Movements created successfully',
      order_id: movements[0].source_id,
    });
  } catch (error) {
    if (error instanceof errors.InsufficientBalanceError) {
      orderTargetStatus = 'failed';
    } else {
      throw error;
    }
  }

  await orderRepository.enqueueOrderToSettle({
    id: movementsToCreate[0].source_id,
    targetStatus: orderTargetStatus,
  });
}

module.exports = {
  createMovement,
};
