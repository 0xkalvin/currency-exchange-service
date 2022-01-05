class InsufficientBalanceError extends Error {
  constructor(message) {
    super(message);
    this.code = 'C_E_INSUFFICIENT_BALANCE';
  }
}

module.exports = {
  InsufficientBalanceError,
};
