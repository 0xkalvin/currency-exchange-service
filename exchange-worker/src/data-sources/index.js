const dynamodb = require('./dynamodb');
const redis = require('./redis');
const sqs = require('./sqs');

module.exports = {
  dynamodb,
  redis,
  sqs,
};
