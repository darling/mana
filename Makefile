.PHONY: dev run build clean test vet fmt lint tidy

BINARY_NAME=mana
BUILD_DIR=bin

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Linker flags
LDFLAGS = -X main.versionStr=$(VERSION) \
          -X main.commitStr=$(COMMIT) \
          -X main.dateStr=$(BUILD_DATE)

dev:
	go run . $(ARGS)

run: dev

build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)
	go clean

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

tidy:
	go mod tidy