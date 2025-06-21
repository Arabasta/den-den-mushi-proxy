GO ?= go
CMD_DIR := cmd
BIN_DIR := bin

# space separated list of binaries
BINARIES := proxy

.PHONY: all build run clean

all: build

build:
	@for b in $(BINARIES); do \
		echo "Building $$b..."; \
		cd $(CMD_DIR)/$$b && $(GO) build -o ../../$(BIN_DIR)/$$b; \
	done

run:
	cd $(CMD_DIR)/$(CMD) && $(GO) run . -config ../../cmd/$(CMD)/config.json

clean:
	rm -rf $(BIN_DIR)
