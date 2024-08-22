# DB_URL=postgresql://root:secret@localhost:5433/defi?sslmode=disable
DB_URL=postgresql://ualxsakaax49i7x90cwk:gv8qRs0tH8ernsr2HaEmXvLMfnQMyx@bbttjw98gs7aswfibtrt-postgresql.services.clever-cloud.com:50013/bbttjw98gs7aswfibtrt

postgres:
	docker run --name defi -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

.PHONY: postgres

createdb:
	docker exec -it defi createdb --username=root --owner=root defi
.PHONY: createdb

dockerlogs:
	docker logs defi
.PHONY: dockerlogs

gencontract:
	solc --abi --bin contract/defi.sol -o build
	abigen --bin=build/CrowdFunding.bin --abi=build/CrowdFunding.abi --pkg=gen --out=gen/crowdFunding.go

dropdb:
	docker exec -it defi dropdb defi
.PHONY: dropdb

createmigrate:
	migrate create -ext sql -dir db/migration -seq schema

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
.PHONY: migrateup

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
.PHONY: migrateup1 migratedown migratedown1 db_docs db_schema

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgress -o docs/schema.sql docs/db.dbml
	
sqlc:
	sqlc generate
.PHONY: sqlc
	
test:
	go test -v -cover ./...

server:
	go run main.go

air:
	air

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/demola234/defiraise/db/sqlc Store
	
.PHONY: mock