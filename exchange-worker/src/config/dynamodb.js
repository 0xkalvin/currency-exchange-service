const {
  DYNAMODB_ENDPOINT,
  DYNAMODB_REGION,
} = process.env;

module.exports = {
  endpoint: DYNAMODB_ENDPOINT,
  region: DYNAMODB_REGION,
};
