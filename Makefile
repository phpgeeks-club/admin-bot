.EXPORT_ALL_VARIABLES:
GOBIN = $(shell pwd)/bin

.PHONY: deps
deps:
	@go mod tidy
	@go mod vendor

.PHONY: mocks
mocks: tools
	@export PATH="$(shell pwd)/bin:$(PATH)"; mockery --config=.mockery.yaml

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test ./...

.PHONY: docker_build
docker_build:
	@docker build -f ./Dockerfile -t geeksonator_dev .

.PHONY: docker_run
docker_run:
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./bin/geeksonator ./cmd/geeksonator

.PHONY: build
build: deps
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./bin/geeksonator ./cmd/geeksonator

.PHONY: tools
tools: deps
	@go install github.com/vektra/mockery/v2@v2.36.0
	@go install github.com/goreleaser/goreleaser@v1.21.2
