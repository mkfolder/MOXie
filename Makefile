compose:
	@docker compose up -d

update:
	@docker compose up -d --build $(NAME)

clean:
	@docker compose down
	@docker image rm moxie-api
	@rm -rf ./volumes/postgres_data
	@mkdir -p ./volumes/postgres_data

recompose: clean compose

.PHONY: compose clean recompose update
