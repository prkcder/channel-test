# Makefile

COMPOSE_FILE ?= docker-compose.yml
SERVICE_NAME ?= scores-api

# Development
build:
	docker compose -f $(COMPOSE_FILE) up --build -d

start:
	docker compose -f $(COMPOSE_FILE) up -d
start-build:
	docker compose -f $(COMPOSE_FILE) up --build -d

stop:
	docker compose -f $(COMPOSE_FILE) down

clean:
	docker compose -f $(COMPOSE_FILE) down -v --remove-orphans

logs:
	docker compose -f $(COMPOSE_FILE) logs -f $(SERVICE_NAME)

restart:
	docker compose -f $(COMPOSE_FILE) down && docker compose -f $(COMPOSE_FILE) up --build -d

ps:
	docker compose -f $(COMPOSE_FILE) ps

shell:
	docker compose -f $(COMPOSE_FILE) exec $(SERVICE_NAME) /bin/sh



# Testing 
test:
	go test ./... -v

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

test-race:
	go test -race ./...

test-store:
	go test -v ./internal/store/...

test-api:
	go test -v ./internal/api/...