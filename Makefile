COMPOSE=docker compose
COMPOSE_BASE=docker-compose.yml
COMPOSE_DEV=docker-compose.dev.yml
COMPOSE_PROD=docker-compose.prod.yml

.PHONY: up down restart logs ps clean

prod:
	$(COMPOSE) --env-file .env -f $(COMPOSE_BASE) -f $(COMPOSE_PROD) up -d --build

dev:
	$(COMPOSE) --env-file .env -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) up --build

ddev:
	$(COMPOSE) --env-file .env -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) up -d --build

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

clean:
	$(COMPOSE) down -v --remove-orphans

auth-migrate-up:
	migrate \
	-path services/auth-service/migrations \
	-database "postgres://app:app@localhost:5432/dark_kitchen?sslmode=disable" \
	up

auth-migrate-down:
	migrate \
	-path services/auth-service/migrations \
	-database "postgres://app:app@localhost:5432/dark_kitchen?sslmode=disable" \
	down 1