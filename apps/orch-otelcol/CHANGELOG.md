<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Customized OpenTelemetry Collector Changelog

## [v0.1.18](https://github.com/open-edge-platform/o11y-charts/tree/apps/orch-otelcol-0.1.18/apps/orch-otelcol) (2025-03-31)

- Initial release
- Container `orch-otelcol` added:
  - Adds [Context Processor](./o11y-otel-contextprocessor/) to allow dynamic setting of `X-Scope-OrgID` header
  - Built with [OpenTelemetry Collector Builder (OCB)](https://opentelemetry.io/docs/collector/custom-collector/) for minimal footprint
