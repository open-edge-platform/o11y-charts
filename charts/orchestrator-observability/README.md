<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# orchestrator-observability

This repository contains definition of Orchestrator Observability stack that handles metrics, logs and traces originating from the Orchestrator itself.

Key components of this observability stack are:

- **Grafana UI** - frontend
- **Grafana Mimir** - metrics backend
- **Grafana Loki** - logs backend
- **Grafana Tempo** - traces backend
- **OpenTelemetry Collector** - pipelines

Grafana Dashboards that are associated with this stack are also kept in this repository.
