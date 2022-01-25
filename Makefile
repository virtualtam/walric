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

race:
	go test -race ./...
.PHONY: race

race10:
	go test -race -count 10 ./...
.PHONY: race10

cover:
	go test -cover ./...
.PHONY: cover

build: $(BUILD_DIR)/walric

$(BUILD_DIR)/%: $(shell find . -type f -name "*.go")
	go build -trimpath -o $@ ./cmd/$*
