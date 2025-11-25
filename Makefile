.PHONY: dev build sync lint format db-up db-down

dev:
	moon :dev

build:
	moon :build

sync:
	moon :sync

lint:
	moon :lint

format:
	moon :format

db-up:
	docker-compose -f docker/db.docker-compose.yaml up -d

db-down:
	docker-compose -f docker/db.docker-compose.yaml down
