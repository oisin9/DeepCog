GO := go
GOFLAGS := -v
BIN_DIR := bin
BIN_NAME := deepcog

.PHONY: all build run clean help

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BIN_NAME) ./cmd/main.go
	@echo "Build successful. Executable is in $(BIN_DIR)/$(BIN_NAME)"
	chmod +x $(BIN_DIR)/$(BIN_NAME)

run: build
	@$(BIN_DIR)/$(BIN_NAME)

clean:
	@rm -rf $(BIN_DIR)
	@$(GO) clean -cache

help:
	@echo "Available Makefile commands:"
	@echo "  make build  - Build project (default)"
	@echo "  make run    - Build and run"
	@echo "  make clean  - Clean build artifacts"
	@echo "  make help   - Show this help message"