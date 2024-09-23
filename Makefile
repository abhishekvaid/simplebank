DB_URL := $(DRIVER_SOURCE)

postgres:
	docker run --name simplebank_pg -p 5432:5432 --network simplebank_network -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:12-alpine

createdb:
	docker exec simplebank_pg createdb --owner=root --username=root simple_bank

dropdb:
	docker exec simplebank_pg dropdb simple_bank

migrateup:
	migrate -path=db/migration -database="$(DB_URL)" up 

migratedown:
	 echo "y" | migrate -path=db/migration -database="$(DB_URL)" down

migrateup1:
	migrate -path=db/migration -database="$(DB_URL)" up 1

migratedown1:
	 echo "y" | migrate -path=db/migration -database="$(DB_URL)" down 1

sqlc: 
	sqlc generate 

test:
 	go test -count=1 ./...

server:
	go run main.go 

dockerrun:
	docker start simplebank_app || docker run --name simplebank_app -p 8081:8081 --network simplebank_network -e DRIVER_SOURCE="postgres://root:secret@simplebank_pg:5432/simple_bank?sslmode=disable" -e SERVER_ADDRESS="0.0.0.0:8081" himavisoft/simplebank:latest

dockerbuild:
	docker rmi himavisoft/simplebank:latest || true && docker build -t himavisoft/simplebank:latest .

mock:
	mockgen --package mockdb --destination ./db/mock/store.go himavisoft.simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 .dockerrun .dockerbuild 