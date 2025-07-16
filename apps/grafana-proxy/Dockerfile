# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Building environment
FROM golang:1.24.5-alpine@sha256:48ee313931980110b5a91bbe04abdf640b9a67ca5dea3a620f01bacf50593396 AS build

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
