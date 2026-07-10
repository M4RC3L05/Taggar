export CGO_ENABLED = 0

GO_FLAGS = -trimpath -ldflags="-w -s"
CURRENT_GIT_TAG := $(shell git describe --tags --exact-match HEAD 2>/dev/null || echo "latest")

.DEFAULT_GOAL: help
.PHONY: help
help:
	@echo "Available targets:"
	@cat $(abspath $(lastword $(MAKEFILE_LIST))) | grep -oP '^[a-zA-Z_-]+(?=:)' | sort | xargs printf "  %s\n"

.PHONY: frontend-install-dep
frontend-install-dep:
	@if [ ! -d "./frontend/node_modules/$(PKG)" ]; then \
		echo "=> Installing $(PKG)" && \
		rm -rf ./frontend/node_modules/$(PKG) && \
		mkdir -p ./frontend/node_modules/$(PKG) && \
		curl --progress-bar -L $(URL) | tar -xz -C ./frontend/node_modules/$(PKG) --strip-components=1; \
	else \
		echo "=> $(PKG) already installed"; \
	fi

.PHONY: frontend-deps
frontend-deps:
	@$(MAKE) -s frontend-install-dep PKG="preact" URL="https://registry.npmjs.org/preact/-/preact-10.29.7.tgz"
	@$(MAKE) -s frontend-install-dep PKG="@digicreon/mucss" URL="https://registry.npmjs.org/@digicreon/mucss/-/mucss-1.4.8.tgz"
	@$(MAKE) -s frontend-install-dep PKG="wouter-preact" URL="https://registry.npmjs.org/wouter-preact/-/wouter-preact-3.10.0.tgz"
	@$(MAKE) -s frontend-install-dep PKG="regexparam" URL="https://registry.npmjs.org/regexparam/-/regexparam-3.0.0.tgz"
	@$(MAKE) -s frontend-install-dep PKG="bootstrap-icons" URL="https://registry.npmjs.org/bootstrap-icons/-/bootstrap-icons-1.13.1.tgz"

.PHONY: frontend-bundle
frontend-bundle: frontend-deps
	rm -rf ./frontend/.dist
	mkdir -p ./frontend/.dist
	cd ./frontend && go run bundle.go && cd ..

.PHONY: code-check
code-check:
	go mod tidy --diff
	golangci-lint run ./...
	golangci-lint fmt --diff-colored ./...
	govulncheck -show verbose -test ./...

.PHONY: main
main: frontend-bundle
	go run $(GO_FLAGS) main.go
