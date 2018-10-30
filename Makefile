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

.PHONY: validate
validate:
	swagger validate ./swagger/api-spec.yml

.PHONY: gen
gen: validate
	swagger generate model \
    		--target=./pkg \
    		--spec=./swagger/api-spec.yml
	swagger generate server \
		--target=./pkg \
		--server-package=api \
		--spec=./swagger/api-spec.yml \
		--exclude-main \
		--name=analyze
		--existing-models=./pkg/models
	cp ./swagger/api-spec.yml ./swagger/ui/api-spec.yml
	statik -f -src=${CURRENT_DIR}/swagger/ui

.PHONY: test
test:
	go test -race ./...
