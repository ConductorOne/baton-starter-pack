GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

ifeq ($(GOOS),windows)
OUTPUT = dist/baton-devolutions.exe
else
OUTPUT = dist/baton-devolutions
endif

.PHONY: build test lint update-deps

build:
	go build -o $(OUTPUT) ./cmd/baton-devolutions

test:
	go test -v ./...

lint:
	golangci-lint run --timeout 3m

update-deps:
	go get -d -u ./...
	go mod tidy

clean:
	rm -rf dist/
