# install_prerequesites:
# 	docker pull quay.io/keycloak/keycloak
# start:
# 	docker run -d -p 8080:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin --name keycloak quay.io/keycloak/keycloak
# stop:
# 	docker-compose down

# restart:
# 	docker-compose down
# 	make start


# Makefile для управления Keycloak и Go-приложением с Docker Compose

# Переменные
DOCKER_COMPOSE = docker-compose
PROJECT_NAME = go_auth

# Задачи
.PHONY: all build up down restart logs keycloak-logs go-logs

# Стандартная команда (build и запуск всех сервисов)
all: build up

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
