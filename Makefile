APP := ghn
DESTINATION_PATH ?= /usr/local/bin/$(APP)
GOLANGCI_LINT_VERSION = "v1.64.0"
GORELEASER_VERSION = "v2.8.0"

.PHONY: golangci-lint-version
golangci-lint-version:
	@echo $(GOLANGCI_LINT_VERSION)

.PHONY: goreleaser-version
goreleaser-version:
	@echo $(GORELEASER_VERSION)

.PHONY: build
build:
	@go build -o $(APP) .
	@go run github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION) build --snapshot

.PHONY: run
run: build
	./$(APP)

## @help:lint:Check the project for linting issues.
.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

## @help:tidy:Update the go.mod
.PHONY: tidy
tidy:
	@go mod tidy

## @help:format:Format the code.
.PHONY: format
format:
	@go fmt ./...

## @help:scan:Scan the code.
.PHONY: scan
scan:
	trufflehog git file://.

## @help:scan:Scan the code.
install: build
	sudo mv ./$(APP) $(DESTINATION_PATH)
