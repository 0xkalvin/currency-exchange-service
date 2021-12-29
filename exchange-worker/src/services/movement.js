const movementRepository = require('../repositories/movement');
const orderRepository = require('../repositories/order');

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

  const movements = await movementRepository.createMovementsAndUpdateBalances(
    movementsToCreate,
    balanceUpdates,
  );

  await orderRepository.enqueueOrderToSettle({
    id: movements[0].source_id,
    targetStatus: 'settled',
  });

  return movements;
}

module.exports = {
  createMovement,
};
