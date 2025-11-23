install:
	set -e
	go mod download
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev:
	air

build:
	go build -o ./bin/bot-service ./cmd/app/main.go

lint:
	golangci-lint run ./...
