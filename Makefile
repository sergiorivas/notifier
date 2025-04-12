.PHONY: build clean install test lint

BINARY_NAME=notify
BUILD_DIR=./build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/notify

clean:
	rm -rf $(BUILD_DIR)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

test:
	go test ./...

lint:
	golangci-lint
