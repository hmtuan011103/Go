.PHONY: all build run test clean lint

# Default
all: build

# Build the application
build:
	go build -o build/bin/server.exe ./cmd/server

# Run the application
run:
	go run ./cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf build/bin/*
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Tidy dependencies
tidy:
	go mod tidy

# Build Docker image
docker-build:
	docker build -f build/package/Dockerfile -t gostructure-app .

# Run with Docker Compose
docker-up:
	docker-compose -f deployments/docker-compose.yaml up -d

# Stop Docker Compose
docker-down:
	docker-compose -f deployments/docker-compose.yaml down
