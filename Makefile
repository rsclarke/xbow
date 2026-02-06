.PHONY: all generate build test clean

all: build

generate:
	go generate ./...

build:
	go build ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...

test-v:
	go test -v ./...

clean:
	go clean ./...
