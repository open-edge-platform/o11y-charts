<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# grafana-proxy

Application `grafana-proxy` is deployed as a sidecar container of Grafana pods
in [orchestrator](../../charts/orchestrator-observability)
and [edgenode](../../charts/edgenode-observability) observability stacks.

The purpose of this container is to set appropriate `X-Scope-OrgID` header with the incoming requests,
based on what the user has access to (from JWT Token).
