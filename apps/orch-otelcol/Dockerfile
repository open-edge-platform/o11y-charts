# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25.5-alpine AS builder

WORKDIR /opt/

COPY . .

RUN apk add --upgrade --no-cache curl=~8 git=~2 make=~4 bash=~5

# Match with the version in builder.yaml
ARG OTEL_VERSION=0.140.0
RUN curl --proto '=https' --tlsv1.2 -fL -o ocb \
    https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${OTEL_VERSION}/ocb_${OTEL_VERSION}_linux_amd64 && \
    chmod +x ocb

RUN ./ocb --config builder.yaml

FROM alpine:3.23@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375

COPY --from=builder /opt/otelcol/otelcontribcol /otelcontribcol

RUN addgroup -S otelcol && adduser -S otelcol -G otelcol
USER otelcol

EXPOSE 4317 55680 55679

ENTRYPOINT ["/otelcontribcol"]
