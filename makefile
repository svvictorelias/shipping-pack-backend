# Makefile - Build and Test Commands
.PHONY: build run test test-ci cover-html test-clean

build:
	go build -o bin/packcalc ./cmd/packcalc

run: build
	set -a; . ./.env.local; set +a; \
	go run ./cmd/packcalc

migrations: build
	set -a; . ./.env.local; set +a; \
	go run ./cmd/migrations

# Run tests with coverage report
test:
	@echo "Running unit tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out | tail -n 1

# Generate HTML coverage report
cover-html:
	@go tool cover -html=coverage.out -o coverage.html
	@echo "coverage.html generated. Open it with your browser."

test-clean:
	@rm -f coverage.out coverage.html
