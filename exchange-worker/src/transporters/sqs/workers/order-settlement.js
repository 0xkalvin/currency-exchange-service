const orderService = require('../../../services/order');

async function settleOrder(message) {
  const payload = JSON.parse(message.Body);

  const {
    customer_id: customerId,
    id,
    target_status: targetStatus,
  } = payload;

  await orderService.settleOrder({
    customerId,
    id,
    targetStatus,
  });
}

module.exports = {
  settleOrder,
};
