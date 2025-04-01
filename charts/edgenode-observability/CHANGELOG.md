<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Edge Node Observability Changelog

## [v0.9.15](https://github.com/open-edge-platform/o11y-charts/tree/charts/edgenode-observability-0.9.15/charts/edgenode-observability) (2025-03-31)

- Initial release
- Chart `edgenode-observability` added:
  - Receives metrics and logs from the Edge Nodes
  - Performs tenancy-aware data processing via [Customized OpenTelemetry Collector](../../apps/orch-otelcol/) pipelines
  - Exposes [Grafana UI](https://grafana.com/oss/grafana/) as frontend
    - Extended with [grafana-proxy](../../apps/grafana-proxy/) for multitenancy support
  - Uses [Grafana Mimir](https://grafana.com/oss/mimir/) (distributed) as metrics backend
  - Uses [Grafana Loki](https://grafana.com/oss/loki/) as logs backend
  - Uses Amazon S3 or S3-compatible storage (MinIO) in `Grafana Mimir` and `Grafana Loki`
  - Scrapes additional Edge Node related metrics from:
    - [Edge Infrastructure Manager Inventory Exporter](https://github.com/open-edge-platform/infra-core/tree/main/exporters-inventory)
    - [App Deployment Manager](https://github.com/open-edge-platform/app-orch-deployment/blob/main/app-deployment-manager)
    - [Observability Tenant Controller](https://github.com/open-edge-platform/o11y-tenant-controller)
