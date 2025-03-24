.DEFAULT_GOAL := build

GO := $(shell which go)
GO_DEFAULT := /usr/local/go/bin/go
ifneq ($(GO),)
else ifneq ($(wildcard $(GO_DEFAULT)),)
    GO := $(GO_DEFAULT)
else
    $(error "Go compiler not found. Please install Go first: https://go.dev/doc/install")
endif
GOFLAGS := -v
BIN_DIR := bin
BIN_NAME := deepcog

.PHONY: all build run clean help

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BIN_NAME) ./cmd/main.go
	@echo "Build successful. Executable is in $(BIN_DIR)/$(BIN_NAME)"

install:
	@echo "=== Starting installation ==="
	@echo "1. Installing executable to /usr/local/bin"
	@sudo install -Dm 755 $(BIN_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)
	
	@echo "2. Creating config directory /etc/deepcog"
	@sudo mkdir -p /etc/deepcog
	
	@echo "3. Installing config file (first install only)"
	@if [ ! -f /etc/deepcog/config.toml ]; then \
		sudo install -Dm 644 config_example.toml /etc/deepcog/config.toml; \
		echo "New config file created: /etc/deepcog/config.toml"; \
	else \
		echo "Preserving existing config file: /etc/deepcog/config.toml"; \
	fi
	
	@echo "4. Updating systemd service"
	@if [ ! -f /etc/systemd/system/deepcog.service ]; then \
		sudo install -Dm 644 deepcog.service /etc/systemd/system/deepcog.service; \
		echo "New systemd service created: /etc/systemd/system/deepcog.service"; \
	else \
		echo "Preserving existing systemd service: /etc/systemd/system/deepcog.service"; \
	fi
	@sudo systemctl daemon-reload
	
	@echo "=== Installation completed ==="
	@echo "Usage instructions:"
	@echo "  Start service: sudo systemctl start deepcog"
	@echo "  Edit config: sudo nano /etc/deepcog/config.toml"
	@echo "  View logs: sudo journalctl -u deepcog -f"

uninstall:
	@rm -f /usr/local/bin/$(BIN_NAME)

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