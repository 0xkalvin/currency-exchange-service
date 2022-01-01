
.PHONY: default all build test down worker

default: infra all

build:
	@docker-compose build --no-cache

infra:
	@docker-compose up -d sqs dynamodb

down:
	@docker-compose down --rmi local --remove-orphans

test:
	@docker-compose run --rm test

all:
	@docker-compose run --rm exchange-worker npm run migrate
	@docker-compose up exchange-worker exchange-api
