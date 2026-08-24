VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo devel)

BUF_VERSION ?= v1.71.0
BUF ?= ./bin/buf

.PHONY: build
# build cli
build:
	mkdir -p ./bin && GOFIPS140=certified go build -ldflags "-w -s -X main.version=$(VERSION)" -o ./bin/ ./cmd/otterscale/...

.PHONY: buf
# install the pinned buf into ./bin
buf:
	@if [ "$$($(BUF) --version 2>/dev/null)" != "$(patsubst v%,%,$(BUF_VERSION))" ]; then \
		echo "Installing buf $(BUF_VERSION)"; \
		GOBIN=$(CURDIR)/bin go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION); \
	fi

.PHONY: proto
# generate Go clients/handlers, the TypeScript package, and the OpenAPI spec from proto
proto: buf
	$(BUF) generate
	$(BUF) generate --template buf.gen.openapi.yaml

.PHONY: proto-lint
# lint proto files
proto-lint: buf
	$(BUF) lint
	$(BUF) format --diff --exit-code

# Override in CI to compare against the pull request base branch, e.g.
# make proto-breaking PROTO_BREAKING_AGAINST="https://github.com/otterscale/otterscale.git#branch=main"
PROTO_BREAKING_AGAINST ?= .git#branch=main

.PHONY: proto-breaking
# check proto files for breaking changes against the main branch
proto-breaking: buf
	$(BUF) breaking --against '$(PROTO_BREAKING_AGAINST)'

.PHONY: vet
# examine code
vet:
	go vet ./...

.PHONY: test
# test code
test:
	go test -coverprofile=coverage.txt ./...

.PHONY: lint
# lint code
lint:
	golangci-lint run

.PHONY: help
# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help