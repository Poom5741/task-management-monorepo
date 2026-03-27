.PHONY: build run test test-coverage clean

APP_NAME=task-management-api
MAIN_PATH=./cmd/server

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...