DB_SOURCE=postgresql://cards:cards@localhost:5432/cards?sslmode=disable

run: 
	@sqlc generate
	@go run cmd/main.go

postgres:
	@echo "Setting up postgres..."
	docker run --name postgres_cards \
		-e POSTGRES_USER=cards  \
		-e POSTGRES_PASSWORD=cards \
		-e POSTGRES_DB=cards \
		-p 5432:5432 -d postgres:alpine

postgresdrop:
	docker rm -f postgres_cards

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

mock: 
	mockgen -destination db/mock/store.go -package mockdb github.com/broemp/red_card/db/sqlc Store
