<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Grafana-proxy

This repository contains code of grafana proxy container used as a sidecar in Grafana pods for
orchestrator and edgenode observability stacks.

The purpose of this container is to set appropriate `X-Scope-OrgID` header with the incoming requests,
based on what the user has access to (from JWT Token).
