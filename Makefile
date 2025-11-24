# Переменные
DOCKER_COMPOSE = docker-compose
SERVICE_REVIEWERS = reviewers
SERVICE_DB = psql
K6_IMAGE = grafana/k6:latest

# Phony цели
.PHONY: up down restart logs logs-db build ps clean dev lint test load-test

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

# Запуск всех тестов
test:
	go test ./... -count=1

# Нагрузочное тестирование (базовое)
load-test:
	@echo "🚀 Запуск нагрузочного теста..."
	k6 run load-test.js

# Комплексное нагрузочное тестирование с разными профилями
load-test-light:
	@echo "🧪 Легкий нагрузочный тест (5 пользователей, 2 минуты)"
	k6 run --vus 5 --duration 2m load-test.js

load-test-medium:
	@echo "⚡ Средний нагрузочный тест (50 пользователей, 5 минут)"
	k6 run --vus 50 --duration 5m load-test.js

load-test-heavy:
	@echo "🔥 Тяжелый нагрузочный тест (200 пользователей, 10 минут)"
	k6 run --vus 200 --duration 10m load-test.js

# Просмотр метрик в реальном времени (требует установленного k6)
k6-dashboard:
	@echo "📈 Запуск k6 с веб-дашбордом..."
	k6 run --out web-dashboard load-test.js

