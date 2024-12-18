.PHONY: all build clean lint format install

BINARY_NAME=gopeek
BUILD_DIR=build
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

install:
	go build -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/main.go

clean:
	rm -rf $(BUILD_DIR)
	rm project_knowledge.md

lint:
	go vet ./...
	go run golang.org/x/lint/golint@latest ./...

format:
	gofmt -s -w $(GO_FILES)