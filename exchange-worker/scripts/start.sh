#!/bin/sh

node /app/src/data-sources/dynamodb/migrations/index.js

exec node /app/src/transporters/sqs/entrypoint.js
