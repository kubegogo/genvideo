.PHONY: help build up down restart logs ps clean

help:
	@echo "GenVideo Development Commands"
	@echo ""
	@echo "  make build    - Build all Docker images"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make restart  - Restart all services"
	@echo "  make logs     - View logs"
	@echo "  make ps       - Show running containers"
	@echo "  make clean    - Remove containers and volumes"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

ps:
	docker-compose ps

clean:
	docker-compose down -v
	rm -rf frontend/node_modules

# Individual services
build-backend:
	docker-compose build backend

build-frontend:
	docker-compose build frontend

logs-backend:
	docker-compose logs -f backend

logs-frontend:
	docker-compose logs -f frontend
