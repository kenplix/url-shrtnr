ifneq (,$(wildcard ./.env))
    include .env
    export
endif

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

.PHONY: install-linter
install-linter:
	@echo "<== Install golangci-lint ==>"
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.46.2

.PHONY: lint
lint: install-linter
	@echo "<== Run golangci-lint ==>"
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: install-linter
	@echo "<== Run golangci-lint fast ==>"
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml

.PHONY: generate
generate:
	@echo "<== Generate files ==>"
	go generate ./...

.PHONY: test
test: generate
	@echo "<== Run unit tests ==>"
	go test -v -race -cover ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	@echo "<== Run go mod tidy ==>"
	go mod tidy -v

.PHONY: run
run: test
	@echo "<== Run application ==>"
	go run cmd/url-shrtnr/main.go
