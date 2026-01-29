# Project settings
APP_NAME := user-management
MAIN := main.go
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

GO := go

# Default target
.PHONY: all
all: build

## Download dependencies
.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod tidy

## Build the binary
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(MAIN)

## Run the application
.PHONY: run
run:
	$(GO) run $(MAIN)

## Run tests
.PHONY: test
test:
	$(GO) test ./...

## Run tests with coverage
.PHONY: test-cover
test-cover:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out

## Lint (requires golangci-lint installed)
.PHONY: lint
lint:
	golangci-lint run

## Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR) coverage.out

## Generate Swagger docs
.PHONY: swagger
swagger:
	swag init
