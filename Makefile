VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo devel)

API_VERSION ?= $(call gomodver,github.com/otterscale/api)
TENANT_OPERATOR_VERSION := v1.0.4
CERT_MANAGER_VERSION    := v1.20.1
FLUX2_VERSION           := v2.8.3

BOOTSTRAP_DIR := manifests/bootstrap
BASE_DIR      := $(BOOTSTRAP_DIR)/base
PLATFORM_DIR  := $(BOOTSTRAP_DIR)/platform

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
	go test -coverprofile=coverage.txt ./...

.PHONY: lint
# lint code
lint:
	golangci-lint run

DOWNLOADED_MANIFESTS := \
	$(BASE_DIR)/cert-manager.yaml \
	$(BASE_DIR)/crds.yaml \
	$(BASE_DIR)/flux2.yaml \
	$(PLATFORM_DIR)/tenant-operator.yaml

.PHONY: bootstrap-manifests
# download bootstrap manifests
bootstrap-manifests: $(DOWNLOADED_MANIFESTS)

$(BASE_DIR)/cert-manager.yaml: | $(BASE_DIR)
	curl -sSL -o $@ \
	  https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml

$(BASE_DIR)/crds.yaml: | $(BASE_DIR)
	curl -sSL -o $@ \
	  https://github.com/otterscale/api/releases/download/$(API_VERSION)/crds.yaml

$(BASE_DIR)/flux2.yaml: | $(BASE_DIR)
	curl -sSL -o $@ \
	  https://github.com/fluxcd/flux2/releases/download/$(FLUX2_VERSION)/install.yaml

$(PLATFORM_DIR)/tenant-operator.yaml: | $(PLATFORM_DIR)
	curl -sSL -o $@ \
	  https://github.com/otterscale/tenant-operator/releases/download/$(TENANT_OPERATOR_VERSION)/install.yaml

$(BASE_DIR) $(PLATFORM_DIR):
	@mkdir -p $@

.PHONY: update-bootstrap
# force re-download all bootstrap manifests
update-bootstrap:
	@rm -f $(DOWNLOADED_MANIFESTS)
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