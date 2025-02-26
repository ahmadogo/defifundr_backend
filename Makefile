# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_URL = ${DB_SOURCE}

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-ps:
	docker-compose ps

docker-build:
	docker-compose build

docker-restart:
	docker-compose restart

# Database commands
postgres:
	docker run --name defi -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it defi createdb --username=root --owner=root defi

dockerlogs:
	docker logs defi

dropdb:
	docker exec -it defi dropdb defi

# Migration commands
createmigrate:
	migrate create -ext sql -dir db/migration -seq schema

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

# Smart contract commands
gencontract:
	solc --abi --bin contract/defi.sol -o build
	abigen --bin=build/CrowdFunding.bin --abi=build/CrowdFunding.abi --pkg=gen --out=gen/crowdFunding.go

# Documentation commands
db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgress -o docs/schema.sql docs/db.dbml

# Development commands
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/api/main.go

air:
	air

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/demola234/defiraise/db/sqlc Store

# Help command
help:
	@echo "Available commands:"
	@echo "  docker-up         - Start all Docker containers"
	@echo "  docker-down       - Stop and remove all Docker containers"
	@echo "  docker-logs       - View logs from all Docker containers"
	@echo "  docker-ps         - List running Docker containers"
	@echo "  docker-build      - Rebuild Docker images"
	@echo "  docker-restart    - Restart all Docker containers"
	@echo ""
	@echo "  postgres          - Run PostgreSQL container"
	@echo "  createdb          - Create database"
	@echo "  dockerlogs        - View PostgreSQL container logs"
	@echo "  dropdb            - Drop database"
	@echo ""
	@echo "  migrateup         - Run all database migrations"
	@echo "  migratedown       - Revert all database migrations"
	@echo "  migrateup1        - Run one database migration"
	@echo "  migratedown1      - Revert one database migration"
	@echo ""
	@echo "  server            - Run the server"
	@echo "  air               - Run the server with hot reload"
	@echo "  test              - Run tests"
	@echo "  sqlc              - Generate SQL code"
	@echo "  mock              - Generate mock code"
	@echo "  gencontract       - Generate contract code"

.PHONY: postgres createdb dockerlogs dropdb createmigrate migrateup migrateup1 migratedown migratedown1 db_docs db_schema sqlc test server air mock gencontract docker-up docker-down docker-logs docker-ps docker-build docker-restart help