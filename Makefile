include .env

run: 
	@sqlc generate
	@go run cmd/main.go

postgres:
	@echo "Setting up postgres..."
	docker run --name postgres_cards \
		-e POSTGRES_USER=$(DB_USER)  \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_DATABASE) \
		-p 5432:5432 -d postgres:alpine

postgresdrop:
	docker rm -f postgres_cards

dropdb:
	@echo "Dropping database..."
	docker exec -it postgres_cards dropdb --username=$(DB_USER) --if-exists 

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_ADDRESS):5432/$(DB_DATABASE)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_ADDRESS):5432/$(DB_DATABASE)?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

mock: 
	mockgen -destination db/mock/store.go -package mockdb github.com/broemp/red_card/db/sqlc Store
