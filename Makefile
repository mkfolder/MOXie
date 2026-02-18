SRC=./cmd/main.go
MIGRATION=./cmd/migrate.go

build:
	@go build -o cynero $(SRC)
	@go build -o migrate $(MIGRATION)

run:
	@go run $(SRC)

migrate:
	@go run $(MIGRATION)

compose:
	@docker compose up -d

recompose:
	@docker compose down
	@docker image rm cynero-cynero
	@rm -rf ./volumes/postgres_data
	@docker compose up -d

update:
	@docker compose up -d --build cynero

.PHONY: build run migrate
