# Include .env file
include .env
export


all: start gen push seed dev


# -------- Development --------
server:
	## Start Air for hot reload
	@APP_ENV=development air

ui:
	## Start Bun development server for UI
	@(cd ./ui && bun run dev)

dev: 
	## Start development server with hot reload
	@$(MAKE) -j ui server

# -------- Services --------
start:
	docker compose --env-file .env -f compose.yaml up -d

stop:
	docker compose --env-file .env -f compose.yaml down

clean:
	docker compose --env-file .env -f compose.yaml down -v

re: clean start

# -------- Misc --------
install:
	@## Install server and ui dependencies
	go mod tidy
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/air-verse/air@latest
	cd ./ui && bun install && cd -



# -------- Database --------
db:
	@## Create a new database
	@goose create $(ARGS) sql

push:

	# wait for the database to be ready
	@until goose up 2> /dev/null; do \
		echo "Waiting for database to be ready..."; \
		sleep 2; \
	done
# 	## Push the database to the latest version
# 	goose up || true
	

down:
	## Migrate the database down by one version
	@goose down 1

migrate:
	## Migrate the database up to the latest version
	@goose $(ARGS)

seed: 
	go run ./cmd/seed/start.go

gen: 
	## Generate Go code from SQL queries
	sqlc generate || true

.PHONY: all server ui dev start stop clean re install db push down migrate seed