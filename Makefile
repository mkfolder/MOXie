compose:
	@docker compose up -d

update:
	@docker compose up -d --build $(NAME)

clean:
	@docker compose down
	@docker image rm moxie-api
	@docker image rm moxie-ui
	@docker image rm moxie-caddy

recompose: clean compose

.PHONY: compose clean recompose update
