.PHONY: build run test lint generate migrate-up migrate-down proto sqlc docker-up docker-down web-install web-dev web-build

APP_NAME := grusha
BUILD_DIR := ./build
MAIN_PATH := ./cmd/grusha

# Build
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

# Testing
test:
	go test -v -race -count=1 ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	golangci-lint run ./...

# Code generation
generate: proto sqlc

proto:
	buf generate

sqlc:
	sqlc generate

# Migrations
migrate-up:
	goose -dir ./migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir ./migrations postgres "$(DATABASE_URL)" down

migrate-create:
	goose -dir ./migrations create $(name) sql

# Docker
docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-build:
	docker compose -f deployments/docker-compose.yml build

# Frontend
web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build
