# Переменные
DOCKER_COMPOSE = docker-compose
PROJECT_NAME = go_auth
SERVICE = auth_service

# Задачи
.PHONY: all build up down restart logs keycloak-logs go-logs

# Стандартная команда (build и запуск всех сервисов)
all: build up

install_prerequesites:
	docker pull quay.io/keycloak/keycloak

# Сборка всех сервисов
build:
	$(DOCKER_COMPOSE) -p $(PROJECT_NAME) build

# Запуск всех сервисов
up:
	$(DOCKER_COMPOSE) -p $(PROJECT_NAME) up -d

# Остановка всех сервисов
down:
	$(DOCKER_COMPOSE) -p $(PROJECT_NAME) down

# Перезапуск всех сервисов
restart: down all

# Просмотр логов всех сервисов
logs:
	$(DOCKER_COMPOSE) -p $(PROJECT_NAME) logs -f

# Логи только Keycloak
keycloak-logs:
	$(DOCKER_COMPOSE) -p $(PROJECT_NAME) logs -f keycloak

# Логи только Go-приложения
go-logs:
	$(DOCKER_COMPOSE) -p $(PROJECT_NAME) logs -f go-auth

tests:
	cd $(SERVICE) && go test ./login -v && go test ./register -v && go test ./logout -v
