.PHONY: all build clean lint format install uninstall test security

BINARY_NAME=gopeek
BUILD_DIR=build
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOLANGCI_LINT_VERSION=v1.62.2


tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gopeek/main.go

install:
	go install ./cmd/gopeek

uninstall:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
	rm -f project_knowledge.md

lint: tools
	golangci-lint run ./...

format:
	gofmt -s -w $(GO_FILES)

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt

security: tools
	go list -json -m all | docker run --rm -i sonatypecommunity/nancy:latest sleuth
	docker run --rm -v $(PWD):/app -w /app securego/gosec:latest -no-fail -fmt=json -out=security-report.json ./...