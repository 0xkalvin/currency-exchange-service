const movementService = require('../../../services/movement');

function createMovement(queueConfig) {
  return async function innerCreateMovement(message) {
    const movementsPayloads = JSON.parse(message.Body);

    await movementService.createMovement(movementsPayloads);
  };
}

module.exports = {
  createMovement,
};
