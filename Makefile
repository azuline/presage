default: build

build:
	go build ./cmd/presage

test:
	go test ./...

lint:
	golangci-lint run --fix --config .golangci.yml ./...
	go mod tidy
