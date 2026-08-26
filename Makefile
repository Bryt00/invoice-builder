# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# You can override this when running make, e.g., make migrate-up DB_DSN="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
DB_DSN ?= "postgres://postgres:raven_db@localhost/invoice-app?sslmode=disable"

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the cmd/web application
.PHONY: run
run:
	go run ./cmd/web

## air: run the application using air for live-reloading
.PHONY: air
air:
	air

## build: build the cmd/web application
.PHONY: build
build: build-css
	go build -o=/tmp/bin/web ./cmd/web

## build-prod: build the cmd/web application for production (linux/amd64 statically linked)
.PHONY: build-prod
build-prod: build-css
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=./bin/invoice-builder ./cmd/web

## build-css: build tailwind css
.PHONY: build-css
build-css:
	npm run build:css

## watch-css: watch tailwind css for changes
.PHONY: watch-css
watch-css:
	npm run watch:css

## test: run all tests
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
