
.PHONY: default all build test down worker

default: infra all

build:
	@docker-compose build --no-cache

all: worker

infra:
	@docker-compose up -d sqs dynamodb

down:
	@docker-compose down --rmi local --remove-orphans

test:
	@docker-compose run --rm test

worker:
	@docker-compose run --rm exchange-worker npm run migrate
	@docker-compose up exchange-worker
