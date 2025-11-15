# Makefile

build:
	docker compose -f docker-compose.yml up --build -d

start:
	docker compose -f docker-compose.yml up -d
start-build:
	docker compose -f docker-compose.yml up --build -d

stop:
	docker compose -f docker-compose.yml down

clean:
	docker compose -f docker-compose.yml down -v --remove-orphans

logs:
	docker compose -f docker-compose.yml logs -f

restart:
	docker compose -f docker-compose.yml down && docker compose -f docker-compose.yml up --build -d
