.PHONY: clean
clean:
	rm -rf ./bin

build-cli: build-migration-cli

build-migration-cli:
	go build -o ./bin/migration ./cmd/migration
