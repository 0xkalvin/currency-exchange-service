const orderService = require('../../../services/order');

async function settleOrder(message) {
  const payload = JSON.parse(message.Body);

  const {
    id,
    target_status: targetStatus,
  } = payload;

  await orderService.settleOrder({
    id,
    targetStatus,
  });
}

module.exports = {
  settleOrder,
};
