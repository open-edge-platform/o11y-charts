<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# edgenode-observability

This stack is intended for handling metrics and logs originating either from the Edge Nodes or from the Edge Orchestrator services that expose data about the Edge Nodes.

Key components of this centralized observability stack are:

- **Grafana UI** - frontend
- **Grafana Mimir** - metrics backend
- **Grafana Loki** - logs backend
- **OpenTelemetry Collector** - pipelines

Grafana Dashboards that are associated with this stack are also kept in [edgenode-dashboards](../edgenode-dashboards).
