VERSION=$(shell git describe --tags --always)
PROTO_FILES=$(shell find api -name *.proto -not -path api/api.proto)
TEST_FILES=$(shell go list ./... | grep -v /api/)

.PHONY: build
# build cli
build:
	mkdir -p ./bin && go build -ldflags "-w -s -X main.version=$(VERSION)" -o ./bin/ ./cmd/otterscale/...

.PHONY: vet
# examine code
vet:
	go vet ./...

.PHONY: test
# test code
test:
	go test -v -coverprofile=coverage.txt $(TEST_FILES)

.PHONY: lint
# lint code
lint:
	golangci-lint run

.PHONY: proto
# generate *.pb.go
proto:
	protoc -I=. -I=third_party -I=third_party/gnostic \
		--go_out=paths=source_relative:. \
		--go_opt=default_api_level=API_OPAQUE \
		--connect-go_out=paths=source_relative:. \
		--connect-go_opt=simple \
		--es_out=web/src/lib \
		--es_opt=target=ts \
		$(PROTO_FILES)

.PHONY: openapi
# generate openapi.yaml
openapi:
	protoc -I=. -I=third_party -I=third_party/gnostic \
		--connect-openapi_out=api \
		--connect-openapi_opt=path=openapi.yaml,short-operation-ids,short-service-tags \
		api/api.proto \
		$(PROTO_FILES)

.PHONY: ceph-deps
# ceph development libraries
ceph-deps:
	apt update && apt install -y libcephfs-dev librbd-dev librados-dev

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