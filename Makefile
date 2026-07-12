export CGO_ENABLED = 0

GO_DIRECT_DEPS := $(shell go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
GO_FLAGS = -trimpath -ldflags="-w -s"
CURRENT_GIT_TAG := $(shell git describe --tags --exact-match HEAD 2>/dev/null || echo "latest")

.DEFAULT_GOAL: help
.PHONY: help
help:
	@echo "Available targets:"
	@cat $(abspath $(lastword $(MAKEFILE_LIST))) | grep -oP '^[a-zA-Z_-]+(?=:)' | sort | xargs printf "  %s\n"

.PHONY:deps-update
deps-update:
	go get -u $(GO_DIRECT_DEPS)
	go mod tidy

.PHONY: code-check
code-check:
	go mod tidy --diff
	golangci-lint run ./...
	golangci-lint fmt --diff-colored ./...
	govulncheck -show verbose -test ./...

.PHONY: main
main:
	go build $(GO_FLAGS) -o ./.bin/main cmd/main.go
