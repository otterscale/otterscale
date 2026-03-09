VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo devel)

API_VERSION ?= $(call gomodver,github.com/otterscale/api)
MODULE_OPERATOR_VERSION := v0.8.3
FLUX2_VERSION := v2.8.1

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
# download bootstrap manifests (FluxCD + module-operator)
bootstrap-manifests: $(BOOTSTRAP_DIR)/crds.yaml $(BOOTSTRAP_DIR)/module-operator.yaml $(BOOTSTRAP_DIR)/flux2.yaml

$(BOOTSTRAP_DIR)/crds.yaml:
	@mkdir -p $(BOOTSTRAP_DIR)
	curl -sSL -o $@ \
	  https://github.com/otterscale/api/releases/download/$(API_VERSION)/crds.yaml

$(BOOTSTRAP_DIR)/module-operator.yaml:
	@mkdir -p $(BOOTSTRAP_DIR)
	curl -sSL -o $@ \
	  https://github.com/otterscale/module-operator/releases/download/$(MODULE_OPERATOR_VERSION)/install.yaml

$(BOOTSTRAP_DIR)/flux2.yaml:
	@mkdir -p $(BOOTSTRAP_DIR)
	curl -sSL -o $@ \
	  https://github.com/fluxcd/flux2/releases/download/$(FLUX2_VERSION)/install.yaml

.PHONY: update-bootstrap
# force re-download all bootstrap manifests
update-bootstrap:
	@rm -f $(BOOTSTRAP_DIR)/flux2.yaml $(BOOTSTRAP_DIR)/module-operator.yaml $(BOOTSTRAP_DIR)/crds.yaml
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

define gomodver
$(shell GOWORK=off go list -m -f '{{if .Replace}}{{.Replace.Version}}{{else}}{{.Version}}{{end}}' $(1) 2>/dev/null)
endef