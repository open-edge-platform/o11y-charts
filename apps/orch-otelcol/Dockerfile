# SPDX-FileCopyrightText: (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.24.2-alpine@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee AS builder

WORKDIR /opt/

COPY . .

RUN apk add --upgrade --no-cache curl=~8 git=~2

# Match with the version in builder.yaml
ARG OTEL_VERSION=0.122.1
RUN curl --proto '=https' --tlsv1.2 -fL -o ocb \
    https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${OTEL_VERSION}/ocb_${OTEL_VERSION}_linux_amd64 && \
    chmod +x ocb

RUN ./ocb --config builder.yaml

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

COPY --from=builder /opt/otelcol/otelcontribcol /otelcontribcol

RUN addgroup -S otelcol && adduser -S otelcol -G otelcol
USER otelcol

EXPOSE 4317 55680 55679

ENTRYPOINT ["/otelcontribcol"]
