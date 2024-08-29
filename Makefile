DB_URL=postgres://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name pg_simplebank -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:12-alpine

createdb:
	docker exec pg_simplebank createdb --owner=root --username=root simple_bank

dropdb:
	docker exec pg_simplebank dropdb simple_bank

migrateup:
	migrate -path=db/migration -database="$(DB_URL)" up 

migratedown:
	 echo "y" | migrate -path=db/migration -database="$(DB_URL)" down

sqlc: 
	sqlc generate 

test:
	go test -v -cover ./... 

server:
	go run main.go 

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server