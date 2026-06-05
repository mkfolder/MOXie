SRC=./cmd/server/main.go
MIGRATION=./cmd/migration/migrate.go

build:
	@go build -o moxie $(SRC)
	@go build -o migrate $(MIGRATION)

run:
	@go run $(SRC)

migrate:
	@go run $(MIGRATION)

compose:
	@docker compose up -d

update:
	@docker compose up -d --build moxie

clean:
	@docker compose down
	@docker image rm moxie-moxie
	@rm -rf ./volumes/postgres_data
	@mkdir -p ./volumes/postgres_data

recompose: clean compose

.PHONY: build run clean recompose migrate
