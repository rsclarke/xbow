.PHONY: all generate build test clean

all: build

generate:
	go generate ./...

build:
	go build ./...

test:
	go test ./...

test-v:
	go test -v ./...

clean:
	go clean ./...
