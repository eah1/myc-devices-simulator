
SCHEMA_DIR := business/db/schema
POSTGRES_URI := postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=utc

start-postgres-test:
	docker run --name postgresTest -e POSTGRES_PASSWORD=postgres -p 5430:5432  -d postgres
	echo "POSTGRES_URI=\"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable\""
	sleep 3

stop-postgres-test:
	docker stop postgresTest
	docker rm postgresTest

goose-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_URI)" goose -dir "$(SCHEMA_DIR)" status

goose-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_URI)" goose -dir "$(SCHEMA_DIR)" up

goose-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_URI)" goose -dir "$(SCHEMA_DIR)" down