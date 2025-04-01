<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Orchestrator Observability Changelog

## [v0.4.12](https://github.com/open-edge-platform/o11y-charts/tree/charts/orchestrator-observability-0.4.12/charts/orchestrator-observability) (2025-03-31)

- Initial release
- Chart `orchestrator-observability` added:
  - Receives metrics, logs and (optionally) traces from the Edge Orchestrator
  - Deploys multiple [OpenTelemetry Collectors](https://opentelemetry.io/docs/collector/):
    - `DaemonSet` (on every node) to collect cluster logs and `kubelet` metrics
    - `Deployment` to collect overall Kubernetes cluster metrics and events
  - Deploys `Prometheus Operator` to collect metrics from Edge Orchestrator services that use `ServiceMonitor` (including [observability-tenant-controller](https://github.com/open-edge-platform/o11y-tenant-controller) for multitenancy project data)
  - Deploys `kube-state-metrics` for collecting metrics from Kubernetes Cluster resources.
  - Exposes [Grafana UI](https://grafana.com/oss/grafana/) as frontend
    - Extended with [grafana-proxy](../../apps/grafana-proxy/) for multitenancy support
  - Uses [Grafana Mimir](https://grafana.com/oss/mimir/) (distributed) as metrics backend
  - Uses [Grafana Loki](https://grafana.com/oss/loki/) as logs backend
  - (Optionally) uses [Grafana Tempo](https://grafana.com/oss/tempo/) as traces backend
  - Uses Amazon S3 or S3-compatible storage (MinIO) in `Grafana Mimir`, `Grafana Loki` and `Grafana Tempo`
