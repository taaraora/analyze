SHELL := /bin/sh

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MAKEFILE_PATH)))


define LINT
	if [ ! -x "`which revive 2>/dev/null`" ]; \
    then \
    	@echo "revive linter not found."; \
    	@echo "Installing linter... into ${GOPATH}/bin"; \
    	go get -u github.com/mgechev/revive ; \
    fi

	@echo "Running code linters..."
	revive
	@echo "Running code linters finished."
endef


.PHONY: default
default: lint


.PHONY: lint
lint:
	@$(call LINT)
