include app.env

# Detect if running on macOS
ifeq ($(shell uname), Darwin)
    DRIVER_SOURCE = postgres://root:secret@localhost:5432/simple_bank?sslmode=disable
endif

# Default target
all:
	@echo "Using DB_URL: $(DB_URL)"

postgres:
	docker run --name pg_simplebank -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres

createdb:
	docker exec pg_simplebank createdb --owner=root --username=root simple_bank

dropdb:
	docker exec pg_simplebank dropdb simple_bank

migrateup:
	migrate -path=db/migration -database="$(DRIVER_SOURCE)" up 

migratedown:
	 echo "y" | migrate -path=db/migration -database="$(DRIVER_SOURCE)" down

migrateup1:
	migrate -path=db/migration -database="$(DRIVER_SOURCE)" up 1

migratedown1:
	 echo "y" | migrate -path=db/migration -database="$(DRIVER_SOURCE)" down 1

sqlc: 
	sqlc generate 

test: 
	go test -count=1 ./...;

server:
	go run main.go 

dockerrun:
	docker start app_simplebank || docker run --name app_simplebank -p 8081:8081 -e DRIVER_SOURCE="$(DRIVER_SOURCE)" himavisoft/simplebank:latest

dockerbuild:
	docker rmi himavisoft/simplebank:latest || true && docker build -t himavisoft/simplebank:latest .

mock:
	mockgen --package mockdb --destination ./db/mock/store.go himavisoft.simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 .dockerrun .dockerbuild 