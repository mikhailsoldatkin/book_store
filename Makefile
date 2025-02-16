include .env
LOCAL_BIN:=$(CURDIR)/bin
REPO:=github.com/mikhailsoldatkin/book_store

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

local-migrations-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${DSN} status -v

local-migrations-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${DSN} up -v

local-migrations-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${DSN} down -v
