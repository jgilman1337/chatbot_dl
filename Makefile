.PHONY: all all_prod clean cli cli_prod discord discord_prod

# Vars
BIN_DIR := bin
GO := go build
GO_PROD := go build -ldflags="-s -w"



all: cli discord
all_prod: cli_prod discord_prod

clean:
	@echo "Cleaning binaries..."
	rm -rf $(BIN_DIR)/*

cli:
	@echo "Building CLI binary..."
	$(GO) -o $(BIN_DIR)/cli ./cmd/cli

cli_prod:
	@echo "Building CLI binary (production version)..."
	$(GO_PROD) -o $(BIN_DIR)/cli_prod ./cmd/cli

discord:
	@echo "Building Discord binary..."
	$(GO) -o $(BIN_DIR)/discord ./cmd/discord

discord_prod:
	@echo "Building Discord binary (production version)..."
	$(GO_PROD) -o $(BIN_DIR)/discord_prod ./cmd/discord

