# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

PROJECT_NAME    := orchestrator-dashboards

.DEFAULT_GOAL   := help
.PHONY: build

# ignore files with Helm template syntax (yamllint can't parse them)
YAMLLINT_IGNORE := ci, trivy, deployments/orchestrator-dashboards/templates/grafana-dashboard-orchestrator/platform-observability.yaml

# Include shared makefile
include ../../common.mk

## CI Mandatory Targets start
dependency-check: ## Unsupported target
	@echo '"make $@" is unsupported'

build: ## Unsupported target
	@echo '"make $@" is unsupported'

lint: lint-markdown lint-yaml lint-json lint-shell lint-helm lint-docker ## Runs all linters

test: ## Unsupported target
	@echo '"make $@" is unsupported'

docker-build: ## Unsupported target
	@echo '"make $@" is unsupported'

helm-build: common-helm-build ## Builds the helm chart

docker-push: ## Unsupported target
	@echo '"make $@" is unsupported'

helm-push: common-helm-push ## Pushes the helm chart
## CI Mandatory Targets end

## Helper Targets start
all: helm-clean helm-build lint ## Runs helm-clean, helm-build, lint targets

clean: helm-clean ## Clean directories

install-tools: ## Installs tools required for the project, requires installed: asdf
	@echo "---MAKEFILE INSTALL-TOOLS---"
	../../scripts/installTools.sh ../.tool-versions
	@echo "---END MAKEFILE INSTALL-TOOLS---"
## Helper Targets end
