<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Grafana Proxy Changelog

## [v0.4.10](https://github.com/open-edge-platform/o11y-charts/tree/apps/grafana-proxy-v0.4.10/apps/grafana-proxy) (2025-03-31)

- Initial release
- Container `grafana-proxy` added:
  - Deployed as a sidecar container with `Grafana UI`
  - Parses user's `JWT Token` forwarded by Datasource (using `oauthPassThru`)
  - Sets appropriate `X-Scope-OrgID` header based on what user has access to
