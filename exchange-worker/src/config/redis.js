const {
  REDIS_HOST,
  REDIS_PORT,
} = process.env;

module.exports = {
  enableOfflineQueue: false,
  host: REDIS_HOST,
  port: Number(REDIS_PORT),
};
