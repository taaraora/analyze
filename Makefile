SHELL := /bin/sh

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MAKEFILE_PATH)))


define LINT
	if [ ! -x "`which golangci-lint 2>/dev/null`" ]; \
    then \
    	@echo "golangci-lint linter not found."; \
    	@echo "Installing linter... into ${GOPATH}/bin"; \
    	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.10.2 ; \
    fi

	@echo "Running code linters..."
	golangci-lint run -j 4 -v
endef


.PHONY: default
default: lint


.PHONY: lint
lint:
	@$(call LINT)
