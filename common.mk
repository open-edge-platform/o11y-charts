# common.mk - common targets for o11y charts repository

# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Makefile Style Guide:
# - Help will be generated from ## comments at end of any target line
# - Use smooth parens $() for variables over curly brackets ${} for consistency
# - When creating targets that run a lint or similar testing tool, print the
#   tool version first so that issues with versions in CI or other remote
#   environments can be caught


# Shell config variable
SHELL                             := bash -eu -o pipefail

## Labels to add Docker/Helm/Service CI meta-data.
LABEL_REVISION                    = $(shell git rev-parse HEAD)
LABEL_CREATED                     ?= $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

VERSION                           ?= $(shell cat VERSION | tr -d '[:space:]')
BUILD_DIR                         ?= ./build

## CHART_NAME is specified in Chart.yaml
CHART_NAME                        ?= $(PROJECT_NAME)
## CHART_VERSION is specified in Chart.yaml
CHART_VERSION                     ?= $(shell grep "version:" ./deployments/$(PROJECT_NAME)/Chart.yaml  | cut -d ':' -f 2 | tr -d '[:space:]')
## CHART_APP_VERSION is modified on every commit
CHART_APP_VERSION                 ?= $(VERSION)
## CHART_BUILD_DIR is given based on repo structure
CHART_BUILD_DIR                   ?= $(BUILD_DIR)/chart/
## CHART_PATH is given based on repo structure
CHART_PATH                        ?= "./deployments/$(CHART_NAME)"
## CHART_RELEASE can be modified here
CHART_RELEASE                     ?= $(PROJECT_NAME)

REGISTRY                          ?= 080137407410.dkr.ecr.us-west-2.amazonaws.com
REGISTRY_NO_AUTH                  ?= edge-orch
REPOSITORY                        ?= o11y
REPOSITORY_NO_AUTH                := $(REGISTRY)/$(REGISTRY_NO_AUTH)/$(REPOSITORY)
DOCKER_IMAGE_NAME                 ?= $(PROJECT_NAME)
DOCKER_IMAGE_TAG                  ?= $(VERSION)

YAMLLINT_IGNORE                   ?= ".github", "ci", "trivy"
YAMLLINT_CONFIG_DATA		      ?= {extends: default, rules: {empty-lines: {max-end: 1}, line-length: {max: 260}, braces: {min-spaces-inside: 0, max-spaces-inside: 1}, brackets: {min-spaces-inside: 0, max-spaces-inside: 1}, document-start: disable}, ignore: [$(YAMLLINT_IGNORE)]}

SH_FILES_TO_LINT                  := $(shell find . \( -type d -name ci \) -prune -o -type f -name '*.sh' -print)

DOCKER_FILES_TO_LINT              := $(shell find . -type f -name 'Dockerfile*' -print )

## Common CI Targets start
common-docker-build: ## Builds docker image
	@echo "---MAKEFILE COMMON-DOCKER-BUILD---"
	docker rmi $(REPOSITORY_NO_AUTH)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) --force
	cp ../../common.mk .
	docker build -f Dockerfile \
		-t $(REPOSITORY_NO_AUTH)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) \
		--build-arg http_proxy="$(http_proxy)" --build-arg https_proxy="$(https_proxy)" --build-arg no_proxy="$(no_proxy)" \
		--platform linux/amd64 --no-cache .
	@rm -rf common.mk
	@echo "---END MAKEFILE COMMON-DOCKER-BUILD---"

common-helm-build: ## Builds the helm chart
	@echo "---MAKEFILE COMMON-HELM-BUILD---"
	yq eval -i '.version = "$(VERSION)"' $(CHART_PATH)/Chart.yaml
	yq eval -i '.appVersion = "$(VERSION)"' $(CHART_PATH)/Chart.yaml
	yq eval -i '.annotations.revision = "$(LABEL_REVISION)"' $(CHART_PATH)/Chart.yaml
	yq eval -i '.annotations.created = "$(LABEL_CREATED)"' $(CHART_PATH)/Chart.yaml
	helm package \
		--app-version=$(CHART_APP_VERSION) \
		--debug \
		--dependency-update \
		--destination $(CHART_BUILD_DIR) \
		$(CHART_PATH)
	@echo "---END MAKEFILE COMMON-HELM-BUILD---"

common-docker-push: ## Pushes the docker image
	@echo "---MAKEFILE COMMON-DOCKER-PUSH---"
	aws ecr create-repository --region us-west-2 --repository-name $(REGISTRY_NO_AUTH)/$(REPOSITORY)/$(DOCKER_IMAGE_NAME) || true
	docker push $(REPOSITORY_NO_AUTH)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@echo "---END MAKEFILE COMMON-DOCKER-PUSH---"

common-helm-push: ## Pushes the helm chart
	@echo "---MAKEFILE COMMON-HELM-PUSH---"
	aws ecr create-repository --region us-west-2 --repository-name $(REGISTRY_NO_AUTH)/$(REPOSITORY)/charts/$(CHART_NAME) || true
	helm push $(CHART_BUILD_DIR)$(CHART_NAME)*.tgz oci://$(REPOSITORY_NO_AUTH)/charts
	@echo "---END MAKEFILE COMMON-HELM-PUSH---"

docker-list: ## Print name of docker container image
	@echo "  $(DOCKER_IMAGE_NAME):"
	@echo "    name: '$(REPOSITORY_NO_AUTH)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)'"
	@echo "    version: '$(DOCKER_IMAGE_TAG)'"
	@echo "    gitTagPrefix: '$(GIT_TAG_PREFIX)'"
	@echo "    buildTarget: '$(PROJECT_NAME)-docker-build'"

helm-list: ## List helm charts, tag format, and versions in YAML format
	@echo "  $(CHART_NAME):" ;\
  echo -n "    "; grep "^version" "${CHART_PATH}/Chart.yaml"  ;\
  echo "    gitTagPrefix: 'charts/$(PROJECT_NAME)-'" ;\
  echo "    outDir: 'charts/$(PROJECT_NAME)/build/chart'" ;\

## Common CI Targets end

## Common helper Targets start
helm-clean: ## Cleans the build directory of the helm chart
	@echo "---MAKEFILE HELM-CLEAN---"
	rm -rf $(CHART_BUILD_DIR)
	@echo "---END MAKEFILE HELM-CLEAN---"

lint-markdown: ## Runs linter for markdown files
	@echo "---MAKEFILE LINT-MARKDOWN---"
	markdownlint-cli2 --config ../../.markdownlint.yml '**/*.md' "!.github" "!ci"
	@echo "---END MAKEFILE LINT-MARKDOWN---"

lint-yaml: ## Runs linter for for yaml files
	@echo "---MAKEFILE LINT-YAML---"
	yamllint -v
	yamllint -f parsable -d '$(YAMLLINT_CONFIG_DATA)' .
	@echo "---END MAKEFILE LINT-YAML---"

lint-json: ## Runs linter for json files
	@echo "---MAKEFILE LINT-JSON---"
	../../scripts/lintJsons.sh "./ci/|./.cache/"
	@echo "---END MAKEFILE LINT-JSON---"

lint-shell: ## Runs linter for shell scripts
	@echo "---MAKEFILE LINT-SHELL---"
	@if [ -n "$(SH_FILES_TO_LINT)" ]; then \
		shellcheck --version; \
		shellcheck $(SH_FILES_TO_LINT); \
	else \
		echo "No shell files to lint."; \
	fi
	@echo "---END MAKEFILE LINT-SHELL---"

lint-helm: ## Runs linter for helm chart
	@echo "---MAKEFILE LINT-HELM---"
	@if [ -f "$(CHART_PATH)" ]; then \
	  helm version; \
	  helm lint --strict $(CHART_PATH); \
	else \
	  echo "CHART_PATH does not exist: $(CHART_PATH)"; \
	fi
	@echo "---END MAKEFILE LINT-HELM---"

lint-docker: ## Runs linter for docker files
	@echo "---MAKEFILE LINT-DOCKER---"
	@if [ -n "$(DOCKER_FILES_TO_LINT)" ]; then \
		shellcheck --version; \
		hadolint $(DOCKER_FILES_TO_LINT); \
	else \
		echo "No Dockerfiles to lint."; \
	fi
	@echo "---END MAKEFILE LINT-DOCKER---"

lint-license: ## Runs license check
	@echo "---MAKEFILE LINT-LICENSE---"
	reuse --version
	reuse --root . lint
	@echo "---END MAKEFILE LINT-LICENSE---"

kind-all: helm-clean common-docker-build kind-load common-helm-build ## Builds docker image and loads it into the kind cluster and builds the helm chart

kind-load: ## Loads docker image into the kind cluster
	@echo "---MAKEFILE KIND-LOAD---"
	kind load docker-image $(REPOSITORY_NO_AUTH)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@echo "---END MAKEFILE KIND-LOAD---"
## Common helper Targets end

## Help Target ####
help: ## Print help for each target
	@echo $(PROJECT_NAME) make targets
	@echo "Target               Makefile:Line    Description"
	@echo "-------------------- ---------------- -----------------------------------------"
	@grep -H -n '^[[:alnum:]_-]*:.* ##' $(MAKEFILE_LIST) \
    | sort -t ":" -k 3 \
    | awk 'BEGIN  {FS=":"}; {sub(".* ## ", "", $$4)}; {printf "%-20s %-20s %s\n", $$3, $$1 ":" $$2, $$4};'
