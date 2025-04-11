BIN_DIR=./bin
CMD_DIR=./cmd
CLI_DIR=./cli

cmd = $(shell ls ${CMD_DIR})
cli = $(shell ls ${CLI_DIR})

# ========================================
# BUILD
# ========================================

.PHONY: build build/cmd build/cli

## build: Builds all cmd and cli packages
build: build/cmd build/cli

## build/cmd: Builds all cmd packages into ./bin
build/cmd: $(cmd)

$(cmd):
	go build -o $(BIN_DIR)/$@ $(CMD_DIR)/$@

## build/cli: Builds cli to ./moon
build/cli: $(cli)

$(cli):
	go build -o moon $(CLI_DIR)

# ========================================
# HELP
# ========================================

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
