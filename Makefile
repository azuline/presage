default: build

build:
	go build ./cmd/presage

lint:
	golangci-lint run --fix --config .golangci.yml ./...
