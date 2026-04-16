# SPDX-FileCopyrightText: (C) 2026 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Building environment
FROM golang:1.26.2-alpine3.23@sha256:c2a1f7b2095d046ae14b286b18413a05bb82c9bca9b25fe7ff5efef0f0826166 AS build

WORKDIR /workspace/apps/grafana-proxy

RUN apk add --upgrade --no-cache make=~4 bash=~5

COPY . .
COPY common.mk /workspace/

RUN make build

FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# Upgrade zlib to fix CVE-2026-22184
RUN apk add --upgrade --no-cache "zlib>=1.3.2-r0"

COPY --from=build /workspace/apps/grafana-proxy/build/grafana-proxy /grafana-proxy

RUN addgroup -S grafanaproxy && adduser -S grafanaproxy -G grafanaproxy
USER grafanaproxy

ENTRYPOINT [ "/grafana-proxy"]
