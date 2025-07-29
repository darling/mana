.PHONY: dev run build clean test vet fmt lint tidy

BINARY_NAME=mana
BUILD_DIR=bin

dev:
	go run . $(ARGS)

run: dev

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

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