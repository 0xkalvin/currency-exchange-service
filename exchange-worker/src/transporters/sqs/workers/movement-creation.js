const movementService = require('../../../services/movement');

async function createMovement(message) {
  const movementsPayloads = JSON.parse(message.Body);

  await movementService.createMovement(movementsPayloads);
}

module.exports = {
  createMovement,
};
