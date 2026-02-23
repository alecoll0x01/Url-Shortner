.PHONY: run test build docker-up docker-down lint

run:
	go run ./cmd/api

test:
	go test -v -race ./internal/...

build:
	go build -ldflags="-w -s" -o bin/api ./cmd/api

docker-up:
	docker compose up --build

docker-down:
	docker compose down

db-up:
	docker compose up db -d

lint:
	golangci-lint run ./...