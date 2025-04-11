BIN_DIR=./bin
CMD_DIR=./cmd
CLI_DIR=./cli
COVERAGE_DIR=./coverage
COVERAGE_OUT=$(COVERAGE_DIR)/cover.out

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
# TIDY
# ========================================

.PHONY: install tidy fmt

## install: Downloads Go dependencies
install:
	@go mod download

## tidy: Cleans unused Go dependencies
tidy:
	@go mod tidy

## fmt: Formats Go code
fmt:
	@gofmt -w .

# ========================================
# CHECK
# ========================================

.PHONY: check check/vet check/fmt

## check: Runs all checks (vet, fmt)
check: check/vet check/fmt

## check/vet: Runs go vet
check/vet:
	@go vet ./...

## check/fmt: Checks for formatting errors. Does not fix
check/fmt:
	@./scripts/fmt-check.sh

# ========================================
# TEST
# ========================================

.PHONY: test test/clean test/coverage

## test: Run all tests with coverage
test: test/clean
	@go test ./... -v -coverprofile=$(COVERAGE_OUT)

## test/clean: Clean coverage directory
test/clean:
	@rm -rf $(COVERAGE_DIR)
	@mkdir $(COVERAGE_DIR)

## test/coverage: Opens test coverage report in browser
test/coverage:
	@go tool cover -html=$(COVERAGE_OUT)

# ========================================
# HELP
# ========================================

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
