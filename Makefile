DOCKER_COMPOSE = docker compose
MANAGER_SERVICE = manager

# Запуск всех сервисов
up:
	$(DOCKER_COMPOSE) up -d

# Остановка всех сервисов
down:
	$(DOCKER_COMPOSE) down

# Перезапуск всех сервисов
restart:
	$(DOCKER_COMPOSE) restart

# Просмотр логов менеджера
logs:
	$(DOCKER_COMPOSE) logs -f $(MANAGER_SERVICE)

# Остановка и удаление всех контейнеров, сетей и образов
clean:
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

# Сборка всех сервисов
build:
	$(DOCKER_COMPOSE) build

# Запуск только менеджера (без воркеров)
up-manager:
	$(DOCKER_COMPOSE) up -d $(MANAGER_SERVICE)

# Остановка только менеджера
down-manager:
	$(DOCKER_COMPOSE) stop $(MANAGER_SERVICE)

# Перезапуск только менеджера
restart-manager:
	$(DOCKER_COMPOSE) restart $(MANAGER_SERVICE)

# Просмотр логов всех сервисов
logs-all:
	$(DOCKER_COMPOSE) logs -f

# Просмотр состояния всех сервисов
ps:
	$(DOCKER_COMPOSE) ps