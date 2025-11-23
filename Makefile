# Переменные
DOCKER_COMPOSE = docker-compose
SERVICE_REVIEWERS = reviewers
SERVICE_DB = psql

# Phony цели
.PHONY: up down restart logs logs-db build ps clean dev lint

# Запуск всех сервисов в фоновом режиме
up:
	$(DOCKER_COMPOSE) up -d

# Остановка всех сервисов
down:
	$(DOCKER_COMPOSE) down

# Перезапуск конкретного сервиса
restart:
	$(DOCKER_COMPOSE) restart $(SERVICE_REVIEWERS)

# Просмотр логов приложения
logs:
	$(DOCKER_COMPOSE) logs -f $(SERVICE_REVIEWERS)

# Просмотр логов БД
logs-db:
	$(DOCKER_COMPOSE) logs -f $(SERVICE_DB)

# Сборка и запуск
build:
	$(DOCKER_COMPOSE) build --no-cache $(SERVICE_REVIEWERS)
	$(DOCKER_COMPOSE) up -d

# Проверка статуса контейнеров
ps:
	$(DOCKER_COMPOSE) ps

# Очистка с удалением томов (осторожно!)
clean:
	$(DOCKER_COMPOSE) down -v

# Запуск в режиме разработки (без фона)
dev:
	$(DOCKER_COMPOSE) up

# Линтинг кода
lint:
	golangci-lint run ./...