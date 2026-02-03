include .env
export

.PHONY: dev build migrate-up migrate-down docker-up docker-down

# Local Development
dev:
	go run cmd/api/main.go

# Docker Compose shortcuts
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

# Database Migrations
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

# Force a migration version (useful if it gets 'dirty')
migrate-force:
	@read -p "Enter version: " v; \
	migrate -path ./migrations -database "$(DATABASE_URL)" force $$v