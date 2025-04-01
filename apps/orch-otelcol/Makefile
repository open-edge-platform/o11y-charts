# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

PROJECT_NAME    := orch-otelcol

.DEFAULT_GOAL   := help
.PHONY: build

YAMLLINT_IGNORE ?= .golangci.yml, ci, trivy

# Include shared makefile
include ../../common.mk

## CI Mandatory Targets start
dependency-check: ## Unsupported target
	@echo '"make $@" is unsupported'

build: ## Unsupported target
	@echo '"make $@" is unsupported'

lint: lint-go lint-markdown lint-yaml lint-json lint-shell lint-docker ## Runs all linters

test: ## Unsupported target
	@echo '"make $@" is unsupported'

docker-build: common-docker-build ## Builds docker image

helm-build: ## Unsupported target
	@echo '"make $@" is unsupported'

docker-push: common-docker-push ## Pushes the docker image

helm-push: ## Unsupported target
	@echo '"make $@" is unsupported'
## CI Mandatory Targets end

## Helper Targets start
all: build lint test ## Runs build ,lint, test targets

clean:
	@echo '"make $@" is unsupported'

lint-go: ## Runs linters for golang source code files
	@echo "---MAKEFILE LINT-GO---"
	cd o11y-otel-contextprocessor/otelcol-dev/contextprocessor && golangci-lint -v run
	@echo "---END MAKEFILE LINT-GO---"

install-tools: ## Installs tools required for the project, requires installed: asdf
	@echo "---MAKEFILE INSTALL-TOOLS---"
	../../scripts/installTools.sh ../.tool-versions
	@echo "---END MAKEFILE INSTALL-TOOLS---"
## Helper Targets end
