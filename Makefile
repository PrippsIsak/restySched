.PHONY: help generate build run test clean install

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies and tools
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest

generate: ## Generate Templ templates
	templ generate

build: generate ## Build the application
	go build -o bin/restysched cmd/server/main.go

run: generate ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test ./... -v

test-coverage: ## Run tests with coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f restysched.db
	find . -name "*_templ.go" -delete

fmt: ## Format code
	go fmt ./...
	templ fmt .

lint: ## Run linter
	go vet ./...

dev: ## Run in development mode with auto-reload (requires air)
	air

.DEFAULT_GOAL := help
