<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# orchestrator-observability

This stack is intended for handling metrics, logs and traces originating from the Edge Orchestrator itself.

Key components of this observability stack are:

- **Grafana UI** - frontend
- **Grafana Mimir** - metrics backend
- **Grafana Loki** - logs backend
- **Grafana Tempo** - traces backend
- **OpenTelemetry Collector** - pipelines

Grafana Dashboards that are associated with this stack are kept in [orchestrator-dashboards](../orchestrator-dashboards).
