# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Building environment
FROM golang:1.26rc2-alpine AS build

WORKDIR /workspace/apps/grafana-proxy

RUN apk add --upgrade --no-cache make=~4 bash=~5

COPY . .
COPY common.mk /workspace/

RUN make build

FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62

COPY --from=build /workspace/apps/grafana-proxy/build/grafana-proxy /grafana-proxy

RUN addgroup -S grafanaproxy && adduser -S grafanaproxy -G grafanaproxy
USER grafanaproxy

ENTRYPOINT [ "/grafana-proxy"]
