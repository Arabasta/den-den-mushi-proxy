GO ?= go
CMD_DIR := cmd
BIN_DIR := bin

# space separated list
BINARIES := proxy control

.PHONY: all build run clean generate check-generate

all: build

build:
	go mod tidy
	mkdir -p $(BIN_DIR)
	@for b in $(BINARIES); do \
		echo "Building $$b..."; \
		cd $(CMD_DIR)/$$b && $(GO) build -o ../../$(BIN_DIR)/$$b; \
		cd - > /dev/null; \
	done

run:
	go mod tidy
	cd $(CMD_DIR)/$(CMD) && $(GO) run . -config ../../cmd/$(CMD)/config.json

clean:
	rm -rf $(BIN_DIR)

#generate:
#	go generate ./...
#
#check-generate:
#	go generate ./...
#	git diff --exit-code || (echo "Generated code out of date"; exit 1)