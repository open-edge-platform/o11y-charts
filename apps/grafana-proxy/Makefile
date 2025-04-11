# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Shell config variable
SHELL := bash -eu -o pipefail

PROJECT_NAME    := grafana-proxy
GIT_TAG_PREFIX  := apps/grafana-proxy-v

VENDOR_DIR      ?= ./vendor

GOCMD           := GOPRIVATE="github.com/open-edge-platform" CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go
GOCMD_TEST      := GOPRIVATE="github.com/open-edge-platform" CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go
GOEXTRAFLAGS    :=-trimpath -gcflags="all=-spectre=all -N -l" -asmflags="all=-spectre=all" -ldflags="all=-s -w -X main.version=$(shell cat ./VERSION)"

.DEFAULT_GOAL   := help
.PHONY: build

YAMLLINT_IGNORE ?= .golangci.yml, vendor, ci, trivy

# Include shared makefile
include ../../common.mk

## CI Mandatory Targets start
dependency-check: ## Unsupported target
	@echo '"make $@" is unsupported'

build: ## Builds grafana-proxy
	@echo "---MAKEFILE BUILD---"
	$(GOCMD) build $(GOEXTRAFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME).go
	@echo "---END MAKEFILE BUILD---"

lint: lint-go lint-markdown lint-yaml lint-json lint-shell lint-docker ## Runs all linters

test: ## Runs tests
	@echo "---MAKEFILE TEST---"
	$(GOCMD_TEST) test ./... --race -coverprofile $(BUILD_DIR)/coverage.out -covermode atomic
	gocover-cobertura < $(BUILD_DIR)/coverage.out > $(BUILD_DIR)/coverage.xml
	@echo "---END MAKEFILE TEST---"

docker-build: vendor common-docker-build ## Builds docker image

helm-build: ## Unsupported target
	@echo '"make $@" is unsupported'

docker-push: common-docker-push ## Pushes the docker image

helm-push: ## Unsupported target
	@echo '"make $@" is unsupported'
## CI Mandatory Targets end

## Helper Targets start
all: clean vendor build lint test ## Runs clean, vendor, build, lint, test targets

clean: ## Deletes build and vendor directories
	@echo "---MAKEFILE CLEAN---"
	rm -rf $(BUILD_DIR)
	rm -rf $(VENDOR_DIR)
	@echo "---END MAKEFILE CLEAN---"

vendor: ## Runs go mod vendor command
	$(GOCMD) mod vendor

lint-go: ## Runs linters for golang source code files
	@echo "---MAKEFILE LINT-GO---"
	golangci-lint -v run
	@echo "---END MAKEFILE LINT-GO---"

install-tools: ## Installs tools required for the project, requires installed: asdf
	@echo "---MAKEFILE INSTALL-TOOLS---"
	../../scripts/installTools.sh ../.tool-versions
	@echo "---END MAKEFILE INSTALL-TOOLS---"
## Helper Targets end
