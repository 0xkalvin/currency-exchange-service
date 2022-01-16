// This file creates dynamodb tables for the development environment
const {
  BatchWriteItemCommand,
  CreateTableCommand,
  DynamoDBClient,
} = require('@aws-sdk/client-dynamodb');

const logger = require('../../../utils/logger')('MIGRATIONS');

const {
  DYNAMODB_ENDPOINT,
  DYNAMODB_REGION,
  BALANCE_TABLE,
  EXCHANGE_TABLE,
  NODE_ENV,
} = process.env;

const dynamoClient = new DynamoDBClient({
  endpoint: DYNAMODB_ENDPOINT,
  region: DYNAMODB_REGION,
});

const migrations = [
  {
    name: 'create-balance-table',
    run: () => {
      const params = {
        TableName: BALANCE_TABLE,
        KeySchema: [
          {
            AttributeName: 'pk',
            KeyType: 'HASH',
          },
          {
            AttributeName: 'sk',
            KeyType: 'RANGE',
          },
        ],
        AttributeDefinitions: [
          {
            AttributeName: 'pk',
            AttributeType: 'S',
          },
          {
            AttributeName: 'sk',
            AttributeType: 'S',
          },
          {
            AttributeName: 'owner_id',
            AttributeType: 'S',
          },
          {
            AttributeName: 'currency_id',
            AttributeType: 'S',
          },
        ],
        GlobalSecondaryIndexes: [
          {
            IndexName: 'owner_id_index',
            KeySchema: [
              {
                AttributeName: 'owner_id',
                KeyType: 'HASH',
              },
              {
                AttributeName: 'sk',
                KeyType: 'RANGE',
              },
            ],
            Projection: {
              ProjectionType: 'ALL',
            },
          },
          {
            IndexName: 'currency_id_index',
            KeySchema: [
              {
                AttributeName: 'currency_id',
                KeyType: 'HASH',
              },
              {
                AttributeName: 'sk',
                KeyType: 'RANGE',
              },
            ],
            Projection: {
              ProjectionType: 'ALL',
            },
          },
        ],
        BillingMode: 'PAY_PER_REQUEST',
        StreamSpecification: {
          StreamEnabled: false,
        },
      };

      if (NODE_ENV === 'production') {
        return Promise.resolve();
      }

      return dynamoClient.send(new CreateTableCommand(params));
    },
  },
  {
    name: 'create-exchange-table',
    run: () => {
      const params = {
        TableName: EXCHANGE_TABLE,
        KeySchema: [
          {
            AttributeName: 'pk',
            KeyType: 'HASH',
          },
          {
            AttributeName: 'sk',
            KeyType: 'RANGE',
          },
        ],
        AttributeDefinitions: [
          {
            AttributeName: 'pk',
            AttributeType: 'S',
          },
          {
            AttributeName: 'sk',
            AttributeType: 'S',
          },
          {
            AttributeName: 'customer_id',
            AttributeType: 'S',
          },
        ],
        GlobalSecondaryIndexes: [
          {
            IndexName: 'customer_id_index',
            KeySchema: [
              {
                AttributeName: 'customer_id',
                KeyType: 'HASH',
              },
              {
                AttributeName: 'sk',
                KeyType: 'RANGE',
              },
            ],
            Projection: {
              ProjectionType: 'ALL',
            },
          },
        ],
        BillingMode: 'PAY_PER_REQUEST',
        StreamSpecification: {
          StreamEnabled: false,
        },
      };

      if (NODE_ENV === 'production') {
        return Promise.resolve();
      }

      return dynamoClient.send(new CreateTableCommand(params));
    },
  },
  {
    name: 'insert-exchange-rates',
    run: async () => {
      const customerId = '123';
      const brlCurrencyId = 'BRL';
      const eurCurrencyId = 'EUR';

      await dynamoClient.send(new BatchWriteItemCommand({
        RequestItems: {
          [EXCHANGE_TABLE]: [
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: brlCurrencyId,
                  },
                  sk: {
                    S: 'EXCHANGE_RATE',
                  },
                  currency_id: {
                    S: 'BRL',
                  },
                  base_currency_id: {
                    S: 'USD',
                  },
                  rate: {
                    N: '17698876',
                  },
                  timestamp: {
                    S: `${new Date().toUTCString()}`,
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: eurCurrencyId,
                  },
                  sk: {
                    S: 'EXCHANGE_RATE',
                  },
                  currency_id: {
                    S: 'EUR',
                  },
                  base_currency_id: {
                    S: 'USD',
                  },
                  rate: {
                    N: '113437680',
                  },
                  timestamp: {
                    S: `${new Date().toUTCString()}`,
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: `${customerId}`,
                  },
                  sk: {
                    S: `${customerId}#ORDER_CUSTOMER#created`,
                  },
                  total: {
                    N: '0',
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: `${customerId}`,
                  },
                  sk: {
                    S: `${customerId}#ORDER_CUSTOMER#settled`,
                  },
                  total: {
                    N: '0',
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: `${customerId}`,
                  },
                  sk: {
                    S: `${customerId}#ORDER_CUSTOMER#failed`,
                  },
                  total: {
                    N: '0',
                  },
                },
              },
            },
          ],
          [BALANCE_TABLE]: [
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: `${customerId}#EUR#available`,
                  },
                  sk: {
                    S: 'BALANCE',
                  },
                  amount: {
                    N: '10000000000',
                  },
                  owner_id: {
                    S: `${customerId}`,
                  },
                  currency_id: {
                    S: 'EUR',
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: `${customerId}#BRL#available`,
                  },
                  sk: {
                    S: 'BALANCE',
                  },
                  amount: {
                    N: '10000000000',
                  },
                  owner_id: {
                    S: `${customerId}`,
                  },
                  currency_id: {
                    S: 'BRL',
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: 'CoolExchange#USD#available',
                  },
                  sk: {
                    S: 'BALANCE',
                  },
                  amount: {
                    N: '10000',
                  },
                  owner_id: {
                    S: 'CoolExchange',
                  },
                  currency_id: {
                    S: 'USD',
                  },
                },
              },
            },
            {
              PutRequest: {
                Item: {
                  pk: {
                    S: `${customerId}`,
                  },
                  sk: {
                    S: `${customerId}#MOVEMENT_OWNER#exchange`,
                  },
                  total: {
                    N: '0',
                  },
                },
              },
            },
          ],
        },
      }));
    },
  },
];

const run = async () => {
  for (const migration of migrations) {
    try {
      logger.info({
        message: `Running ${migration.name}`,
      });

      await migration.run();
    } catch (error) {
      if (error.name === 'ResourceInUseException') {
        logger.info({
          message: `Already exists ${migration.name}`,
        });

        break;
      } else {
        logger.fatal({
          message: `Failed to run ${migration.name}`,
          error_message: error.message,
          error_stack: error.stack,
        });

        throw error;
      }
    } finally {
      logger.info({
        message: `Done ${migration.name}`,
      });
    }
  }
};

run();
