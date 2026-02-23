VERSION=$(shell git describe --tags --always)

FLUX2_VERSION := v2.4.0
OTTERSCALE_OPERATOR_VERSION := v0.2.6

BOOTSTRAP_DIR := manifests/bootstrap

.PHONY: build
# build cli
build: bootstrap-manifests
	mkdir -p ./bin && GOFIPS140=latest go build -ldflags "-w -s -X main.version=$(VERSION)" -o ./bin/ ./cmd/otterscale/...

.PHONY: vet
# examine code
vet:
	go vet ./...

.PHONY: test
# test code
test:
	go test -v -coverprofile=coverage.txt ./...

.PHONY: lint
# lint code
lint:
	golangci-lint run

.PHONY: bootstrap-manifests
# download bootstrap manifests (FluxCD + otterscale-operator)
bootstrap-manifests: $(BOOTSTRAP_DIR)/flux2.yaml $(BOOTSTRAP_DIR)/otterscale-operator.yaml

$(BOOTSTRAP_DIR)/flux2.yaml:
	@mkdir -p $(BOOTSTRAP_DIR)
	curl -sSL -o $@ \
	  https://github.com/fluxcd/flux2/releases/download/$(FLUX2_VERSION)/install.yaml

$(BOOTSTRAP_DIR)/otterscale-operator.yaml:
	@mkdir -p $(BOOTSTRAP_DIR)
	curl -sSL -o $@ \
	  https://github.com/otterscale/otterscale-operator/releases/download/$(OTTERSCALE_OPERATOR_VERSION)/install.yaml

.PHONY: update-bootstrap
# force re-download all bootstrap manifests
update-bootstrap:
	@rm -f $(BOOTSTRAP_DIR)/flux2.yaml $(BOOTSTRAP_DIR)/otterscale-operator.yaml
	$(MAKE) bootstrap-manifests

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