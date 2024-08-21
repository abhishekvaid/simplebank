DB_URL=postgres://root:secret@localhost:54321/simple_bank?sslmode=disable

postgres:
	docker run --name pg_simplebank --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine 

createdb:
	docker exec -it pg_simplebank createdb --owner=root --username=root simple_bank 

dropdb:
	docker exec -it pg_simplebank dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "$(DB_URL)"  up 

migratedown:
	migrate --path db/migrations -database "$(DB_URL)" down 

.PHONY: postgres createdb dropdb migrateup migratedown
