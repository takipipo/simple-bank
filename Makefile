postgres-start:
	@docker run --name postgres12 -d -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:12-alpine
postgres-stop:
	@docker rm --force postgres12
createdb:
	@docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	@docker exec -it postgres12 dropdb simple_bank

migrate-up:
	@migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate-down:
	@migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc-generate:
	@sqlc generate
.PHONY: createdb postgres-start dropdb postgres-stop migrate-up sqlc-generate