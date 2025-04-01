<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Edge Node Dashboards Changelog

## [v0.2.29](https://github.com/open-edge-platform/o11y-charts/tree/charts/edgenode-dashboards-0.2.29/charts/edgenode-dashboards) (2025-03-31)

- Initial release
- Chart `edgenode-dashboards` added:
  - Provides sets of dashboards for Grafana in [edgenode-observability](../edgenode-observability/)
  - Adds Host dashboards `ConfigMap` that contains:
    - Host Performance
    - Agent Logs
    - System Logs
    - Provisioning Logs
    - Healthcheck
  - Adds Cluster dashboards `ConfigMap` that contains:
    - Cluster
    - Pod Logs
    - Edge Node Pods
    - Virtual Machines
