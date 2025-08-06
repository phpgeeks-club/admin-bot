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
	@docker build --no-cache -f ./Dockerfile -t geeksonator_dev .

.PHONY: docker_run
docker_run:
	@docker run -d --env-file=./.env --name geeksonator.dev geeksonator_dev .

.PHONY: tools
tools: deps
	@go install github.com/vektra/mockery/v2@v2.36.0
