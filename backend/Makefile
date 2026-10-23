.PHONY: build run dev clean test deps

# Variables
BINARY_NAME=paypilot-api
BUILD_DIR=./bin

# Build the application
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run: build
	@echo "Running..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@air

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f paypilot.db

# Run tests
test:
	@go test -v ./...

# Download dependencies
deps:
	@go mod download
	@go mod tidy

# Install dev dependencies
dev-deps:
	@go install github.com/air-verse/air@latest

# Create .env from example
env:
	@cp .env.example .env
	@echo "Created .env file. Please update with your credentials."
