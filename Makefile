# cmdX Makefile — developer convenience targets

BINARY   := cmdx
CMD_DIR  := ./cmd
TEST_FLAGS := -v -race

.PHONY: all build test vet lint deps clean install help

## all: build and test
all: deps build test

## deps: install/update all Go dependencies
deps:
	go mod tidy
	go mod download

## build: compile the cmdx binary
build:
	go build -o $(BINARY).exe $(CMD_DIR)/

## test: run the full test suite
test:
	go test $(TEST_FLAGS) ./...

## vet: run go vet
vet:
	go vet ./...

## lint: run golangci-lint (must be installed: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

## clean: remove build artifacts
clean:
	rm -f $(BINARY).exe $(BINARY)

## install: build and copy binary to PATH (Linux/macOS)
install: build
	cp $(BINARY) $(HOME)/.local/bin/$(BINARY)
	@echo "Installed to $(HOME)/.local/bin/$(BINARY)"

## help: show this help
help:
	@grep -E '^##' Makefile | sed 's/## /  /'
