.PHONY: install dev build lint format

# Load environment variables from ./config/.env
-include config/.env
export

# dsn constructed from environment variables
db_url=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

install:
	set -e; \
	go mod download; \
	go install github.com/air-verse/air@latest; \
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev:
	air

build:
	go build -o ./bin/bot-service ./cmd/app/main.go

lint:
	golangci-lint run

format:
	go fmt ./...

migrate-create:
	docker run -i -v ./database/migrations:/database/migrations --network host migrate/migrate create -ext sql -dir /database/migrations $(name)

migrate-up:
	docker run -i -v ./database/migrations:/database/migrations --network host migrate/migrate -path /database/migrations -database $(db_url) -verbose up $(steps)

migrate-down:
	docker run -i -v ./database/migrations:/database/migrations --network host migrate/migrate -path /database/migrations -database $(db_url) -verbose down $(steps)

migrate-force:
	docker run -i -v ./database/migrations:/database/migrations --network host migrate/migrate -path /database/migrations -database $(db_url) -verbose force $(version)
