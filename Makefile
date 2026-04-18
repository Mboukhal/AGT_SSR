

# Include .env file
include .env
export


all: dev



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
	## Create a new database
	@goose -dir ${GOOSE_MIGDIR} create $(ARGS) sql

push:
	## Push the database to the latest version
	@mkdir -p ${DATABASE_DIR}; 
	@goose -dir ${GOOSE_MIGDIR} sqlite3 $${DATABASE_URL} up || true
	@go run ./seed/db.go

down:
	## Migrate the database down by one version
	@goose -dir ${GOOSE_MIGDIR} sqlite3 $${DATABASE_URL} down 1

migrate:
	## Migrate the database up to the latest version
	@goose -dir ${GOOSE_MIGDIR} sqlite3 $${DATABASE_URL} $(ARGS)

seed: push gen
	go run ./cmd/seed/db.go

.PHONY: all server ui dev start stop clean re install