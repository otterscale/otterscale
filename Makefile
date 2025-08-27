VERSION=$(shell git describe --tags --always)
PROTO_FILES=$(shell find api -name *.proto)

.PHONY: build
# build cli
build:
	mkdir -p ./bin && go build -ldflags "-w -s -X main.version=$(VERSION)" -o ./bin/ ./cmd/otterscale/...

.PHONY: vet
# examine code
go-vet:
	go vet ./...

.PHONY: test
# test code
test:
	go test -race ./...
	mkdir -p coverage/unit
	go test -cover ./... -args -test.gocoverdir="$$PWD/coverage/unit"
	go tool covdata textfmt -i=./coverage/unit -o coverage/profile
	go tool cover -func coverage/profile

.PHONY: lint
# lint code
lint:
	golangci-lint run

.PHONY: proto
# generate *.pb.go
proto:
	protoc -I=. \
		--go_out=paths=source_relative:. \
		--go_opt=default_api_level=API_OPAQUE \
		--connect-go_out=paths=source_relative:. \
		--es_out=web/src/lib \
		--es_opt=target=ts \
		$(PROTO_FILES)

.PHONY: service-image
# build backend image
service-image:
	docker build -f Dockerfile --network host -t otterscale:$(VERSION) .

.PHONY: web-image
# build frontend image
web-image:
	cd web && docker build -f Dockerfile --network host -t otterscale-web:$(VERSION) .

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