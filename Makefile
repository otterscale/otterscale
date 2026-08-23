VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo devel)

.PHONY: build
# build cli
build:
	mkdir -p ./bin && GOFIPS140=certified go build -ldflags "-w -s -X main.version=$(VERSION)" -o ./bin/ ./cmd/otterscale/...

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