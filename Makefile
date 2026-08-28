# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Database connection DSN (override via make migrate-up DB_DSN="...")
DB_DSN ?= "postgres://postgres:raven_db@localhost/invoice-app?sslmode=disable"

# ==================================================================================== #
# DEVELOPMENT & BUILD
# ==================================================================================== #

## run: run the Go API backend (cmd/api)
.PHONY: run
run: run-api

## run-api: run the Go API backend (cmd/api)
.PHONY: run-api
run-api:
	go run ./cmd/api

## run-frontend: run the Vue 3 Vite development server
.PHONY: run-frontend
run-frontend:
	cd frontend && npm run dev

## air: run the API application using air for live-reloading
.PHONY: air
air:
	air

## build-api: build the Go API binary
.PHONY: build-api
build-api:
	go build -o=/tmp/bin/api ./cmd/api

## build-prod: build the Go API application for production (linux/amd64 statically linked)
.PHONY: build-prod
build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=./bin/invoice-builder ./cmd/api

## build-frontend: build the Vue 3 SPA production bundle
.PHONY: build-frontend
build-frontend:
	cd frontend && npm run build

## build-all: build both Go API backend and Vue 3 frontend
.PHONY: build-all
build-all: build-api build-frontend

## test: run all Go unit tests
.PHONY: test
test:
	go test -v -race ./...

## create-admin: create or promote an admin user (e.g., make create-admin email=admin@example.com password=secret name="Admin User")
.PHONY: create-admin
create-admin:
	go run ./cmd/create-admin -email="$(email)" -password="$(password)" -name="$(name)"


# ==================================================================================== #
# DATABASE MIGRATIONS (using goose)
# ==================================================================================== #

## migrate-status: check the status of all migrations
.PHONY: migrate-status
migrate-status:
	goose -dir migrations postgres $(DB_DSN) status

## migrate-up: apply all up migrations
.PHONY: migrate-up
migrate-up:
	goose -dir migrations postgres $(DB_DSN) up

## migrate-down: apply all down migrations (rollback one step)
.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres $(DB_DSN) down

## migrate-reset: rollback all migrations
.PHONY: migrate-reset
migrate-reset:
	goose -dir migrations postgres $(DB_DSN) down-to 0

## migrate-create name=$1: create a new migration (e.g., make migrate-create name=add_users)
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=my_migration"; \
		exit 1; \
	fi
	goose -dir migrations create $(name) sql

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
