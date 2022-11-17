# Easily setup the development environment by running make <command name>
postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Madara123 -d postgres:15-alpine
createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres15 dropdbdb --username=root --owner=root simple_bank
migrateup:
	migrate -path db/migrations -database "postgresql://root:Madara123@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:Madara123@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate

.PHONY: postgres dropdb createdb migrateup migratedown