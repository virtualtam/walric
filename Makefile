BUILD_DIR ?= build

all: lint test build
.PHONY: all

clean:
	go clean -i ./...
	rm -rf $(BUILD_DIR)

lint:
	golangci-lint run ./...
.PHONY: lint

test:
	go test ./...
.PHONY: test

build: $(BUILD_DIR)/redwall

$(BUILD_DIR)/%: $(shell find . -type f -name "*.go")
	go build -trimpath -o $@ ./cmd/$*
