SHELL := /bin/bash

.PHONY: help
help: ## Display this help text
	@echo 'Perform common development tasks.'
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | awk -F ':.*?## ' 'NF==2 {printf "\033[36m  %-25s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Runs linting
	./build/run.sh "./build/lint.sh"

.PHONY: test
test: ## Run all tests
	./build/run.sh "./build/test.sh"