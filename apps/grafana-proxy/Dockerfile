# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Building environment
FROM golang:1.24.6-alpine@sha256:c8c5f95d64aa79b6547f3b626eb84b16a7ce18a139e3e9ca19a8c078b85ba80d AS build

WORKDIR /workspace/apps/grafana-proxy

RUN apk add --upgrade --no-cache make=~4 bash=~5

COPY . .
COPY common.mk /workspace/

RUN make build

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

COPY --from=build /workspace/apps/grafana-proxy/build/grafana-proxy /grafana-proxy

RUN addgroup -S grafanaproxy && adduser -S grafanaproxy -G grafanaproxy
USER grafanaproxy

ENTRYPOINT [ "/grafana-proxy"]
