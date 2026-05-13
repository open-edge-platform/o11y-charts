# SPDX-FileCopyrightText: (C) 2026 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.26.3-alpine3.23@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS builder

WORKDIR /opt/

COPY . .

RUN apk add --upgrade --no-cache curl=~8 git=~2 make=~4 bash=~5

# Match with the version in builder.yaml
ARG OTEL_VERSION=0.152.0
RUN curl --proto '=https' --tlsv1.2 -fL -o ocb \
    https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${OTEL_VERSION}/ocb_${OTEL_VERSION}_linux_amd64 && \
    chmod +x ocb

RUN ./ocb --config builder.yaml

FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# Upgrade zlib to fix CVE-2026-22184
RUN apk add --upgrade --no-cache "zlib>=1.3.2-r0" "libcrypto3>=3.5.6-r0" "libssl3>=3.5.6-r0"

COPY --from=builder /opt/otelcol/otelcontribcol /otelcontribcol

RUN addgroup -S otelcol && adduser -S otelcol -G otelcol
USER otelcol

EXPOSE 4317 55680 55679

ENTRYPOINT ["/otelcontribcol"]
