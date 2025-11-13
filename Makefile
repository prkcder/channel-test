# Makefile

start-build:
	docker compose -f docker-compose.yml up --build

start-up:
	docker compose -f docker-compose.yml up -d
start:
	docker compose -f docker-compose.yml up --build -d

stop:
	docker compose -f docker-compose.yml down

stop-clean:
	docker compose -f docker-compose.yml down -v

logs:
	docker compose -f docker-compose.yml logs -f

restart:
	docker compose -f docker-compose.yml down && docker compose -f docker-compose.yml up --build
