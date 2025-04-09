.PHONY: clean
clean:
	rm -rf ./bin

build: build-tools

build-tools: build-migration

build-migration:
	go build -o ./bin/migration ./tools/migration

test:
	go test -cover -coverprofile=coverage.out ./... -v

coverage:
	go tool cover -html=coverage.out
