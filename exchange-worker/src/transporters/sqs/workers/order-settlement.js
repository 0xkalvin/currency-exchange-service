const orderService = require('../../../services/order');

function settleOrder(queueConfig) {
  return async function innerSettleOrder(message) {
    const payload = JSON.parse(message.Body);

    const {
      id,
      target_status: targetStatus,
    } = payload;

    await orderService.settleOrder({
      id,
      targetStatus,
    });
  };
}

module.exports = {
  settleOrder,
};
