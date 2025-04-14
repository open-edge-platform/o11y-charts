# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Shell config variable
SHELL := bash -eu -o pipefail

MAKEFLAGS += --no-print-directory

SUBPROJECTS := apps/grafana-proxy apps/orch-otelcol charts/edgenode-dashboards charts/edgenode-observability charts/orchestrator-dashboards charts/orchestrator-observability

.DEFAULT_GOAL := help
.PHONY: all build lint test clean help

all: build lint test
	@# Help: Runs build, lint, test targets for all subprojects

build:
	@# Help: Runs build target in all subprojects
	@echo "---MAIN MAKEFILE BUILD---"
	for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir build; done
	@echo "---END MAIN MAKEFILE BUILD---"

lint:
	@# Help: Runs lint target in all subprojects
	@echo "---MAIN MAKEFILE LINT---"
	@for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir lint; done
	@echo "---END MAIN MAKEFILE LINT---"

test:
	@# Help: Runs test target in all subprojects
	@echo "---MAIN MAKEFILE TEST---"
	for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir test; done
	@echo "---END MAIN MAKEFILE TEST---"

docker-build:
	@# Help: Runs docker-build target in all subprojects
	@echo "---MAIN MAKEFILE DOCKER-BUILD---"
	for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir docker-build; done
	@echo "---END MAIN MAKEFILE DOCKER-BUILD---"

APPS = apps/grafana-proxy apps/orch-otelcol
docker-list: ## list all docker containers
	@echo "images:"
	@for dir in $(APPS); do $(MAKE) -C $$dir $@; done

helm-build:
	@# Help: Runs helm-build target in all subprojects
	@echo "---MAIN MAKEFILE HELM-BUILD---"
	for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir helm-build; done
	@echo "---END MAIN MAKEFILE HELM-BUILD---"

CHARTS = charts/edgenode-dashboards charts/edgenode-observability charts/orchestrator-dashboards charts/orchestrator-observability
helm-list: ## list all docker containers
	@echo "charts:"
	@for dir in $(CHARTS); do $(MAKE) -C $$dir $@; done

clean:
	@# Help: Runs clean target in all subprojects
	@echo "---MAIN MAKEFILE CLEAN---"
	for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir clean; done
	@echo "---END MAIN MAKEFILE CLEAN---"

install-tools:
	@# Help: Runs install-tools target in all subprojects
	@echo "---MAIN MAKEFILE INSTALL-TOOLS---"
	for dir in $(SUBPROJECTS); do $(MAKE) -C $$dir install-tools; done
	@echo "---END MAIN MAKEFILE INSTALL-TOOLS---"

grafana-proxy-%:
	@# Help: Runs grafana-proxy subproject's tasks, e.g. grafana-proxy-lint
	$(MAKE) -C apps/grafana-proxy $*

orch-otelcol-%:
	@# Help: Runs orch-otelcol subproject's tasks, e.g. orch-otelcol-lint
	$(MAKE) -C apps/orch-otelcol $*

edge-dash-%:
	@# Help: Runs edgenode-dashboards subproject's tasks, e.g. edgenode-dashboards-lint
	$(MAKE) -C charts/edgenode-dashboards $*

edge-obs-%:
	@# Help: Runs edgenode-observability subproject's tasks, e.g. edgenode-observability-lint
	$(MAKE) -C charts/edgenode-observability $*

orch-dash-%:
	@# Help: Runs orchestrator-dashboards subproject's tasks, e.g. orchestrator-dashboards-lint
	$(MAKE) -C charts/orchestrator-dashboards $*

orch-obs-%:
	@# Help: Runs orchestrator-observability subproject's tasks, e.g. orchestrator-observability-lint
	$(MAKE) -C charts/orchestrator-observability $*

help:
	@printf "%-20s %s\n" "Target" "Description"
	@printf "%-20s %s\n" "------" "-----------"
	@grep -E '^[a-zA-Z0-9_%-]+:|^[[:space:]]+@# Help:' Makefile | \
	awk '\
		/^[a-zA-Z0-9_%-]+:/ { \
			target = $$1; \
			sub(":", "", target); \
		} \
		/^[[:space:]]+@# Help:/ { \
			if (target != "") { \
				help_line = $$0; \
				sub("^[[:space:]]+@# Help: ", "", help_line); \
				printf "%-20s %s\n", target, help_line; \
				target = ""; \
			} \
		}' | sort -k1,1
