.PHONY: clean
clean:
	rm -rf ./bin

build-cli: build-migration-cli

build-migration-cli:
	go build -o ./bin/migration ./cmd/migration

test:
	go test -cover -coverprofile=coverage.out ./... -v

coverage:
	go tool cover -html=coverage.out
