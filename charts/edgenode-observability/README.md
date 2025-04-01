<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# edgenode-observability

This repository contains definition of Edge Node Observability stack that handles metrics and logs either originating from the Edge Nodes or exposed by the Orchestrator but related to Edge Nodes.

Key components of this centralized observability stack are:

- **Grafana UI** - frontend
- **Grafana Mimir** - metrics backend
- **Grafana Loki** - logs backend
- **OpenTelemetry Collector** - pipelines

Grafana Dashboards that are associated with this stack are also kept in this repository.
