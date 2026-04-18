




all: start


# dev:


start:
	docker compose --env-file .env -f compose.yaml up -d

stop:
	docker compose --env-file .env -f compose.yaml down

clean:
	docker compose --env-file .env -f compose.yaml down -v

re: clean start

.PHONY: all start stop clean re