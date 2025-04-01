<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Origin

Current directory was copied from <https://github.com/springernature/o11y-otel-contextprocessor/commit/cf30922e0639164a7d10bd68d34b3117447f6894> and adapted to meet specific needs.

## Context OTEL processor

A processor to manipulate the Context Medatada in OTEL Collector. This processor allows to move resource attributes of Metrics, Logs and Traces to the context, which makes possible to use extensions like `headers_setter` dynamically.

Typical use cases:

- Be able to dynamically define tenants for Mimir/Cortex, Loki and Tempo
- Dynamically define metadata attributes in the context, to offer a link to pass resource attribute to extensions
- Change metadata generated from the receivers

The purpose of this code is quickly integrate a new component "Context processor" in our OpenTelemetry collectors. This is a temporary solution while this functionality is not available in upstream.

You can read about what Context Processor does in [otelcol-dev/contextprocessor/README.md](./otelcol-dev/contextprocessor/README.md). And you can see an example configuration in [otelcol-dev/otelcol.yaml](./otelcol-dev/otelcol.yaml)

## About

To get starting contributing in OpenTelemetry, please have a look at these resources:

- OpenTelemetry Collector: <https://github.com/open-telemetry/opentelemetry-collector/blob/main/CONTRIBUTING.md>
- OpenTelemetry Collector-Contrib: <https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/CONTRIBUTING.md>

## Building the collector

1. Get golang installed: <https://go.dev/doc/install>
1. Install the [ocb](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder) tool: <https://github.com/open-telemetry/opentelemetry-collector/releases>

If you want to get to your own custom collector compiled or updated:

1. Update the file `builder-config.yaml`, specially the new custom version in `version`. To completely update all components and add new ones as in the [Contrib distribution](https://github.com/open-telemetry/opentelemetry-collector-contrib), you can copy it from upstream <https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/cmd/otelcontribcol/builder-config.yaml> and adjust `name`, `output_path`, `otelcol_version` and all the gomod versions and make sure you remove the section `replaces`.

    ```yaml
    dist:
        name: otelcol-dev
        description: Local OpenTelemetry Collector Contrib binary, testing only.
        output_path: ./otelcol-dev
        version: 0.105.0-sn3
        otelcol_version: 0.105.0
   ```

1. Add the Context processor component to the processors section of `builder-config.yaml` (with the same reference you used in the `go mod init ...`). Create another section to replace this module name with the local directory where the code is (`./contextprocessor` in this case). Put the same version as the rest of the modules (this only to avoid a build error with `ocb`, is not important for the compiler!)

   ```yaml
   processors:
   - gomod: github.com/springernature/o11y-otel-contextprocessor/processor/contextprocessor v0.105.0

   replaces:
   - github.com/springernature/o11y-otel-contextprocessor/processor/contextprocessor => ./contextprocessor
   ```

1. Run `ocb --config builder-config.yaml`. A new directory `otelcol-dev` (or what you have defined in `output_path`) is created. The code generated should not be touched, all changes should be done via the `builder-config.yaml` file.
1. Go to new `otelcold-dev` and you will see a new executable `otelcol-dev` (or the value defined in `name`). You can execute it by providing a basic configuration file `otelcold-dev --config otelcol.yaml`:

    ```shell
    cat <<EOF > otelcol.yaml
    extensions:
    health_check:
    pprof:
        endpoint: 0.0.0.0:1777
    zpages:
        endpoint: 0.0.0.0:55679

    receivers:
    prometheus:
        config:
        scrape_configs:
            - job_name: 'otel-collector'
            scrape_interval: 30s
            static_configs:
                - targets: ['0.0.0.0:8888']

    exporters:
    debug:
        verbosity: detailed

    service:
    telemetry:
        logs:
        level: "debug"
    pipelines:
        metrics:
        receivers: [prometheus]
        processors: []
        exporters: [debug]
    extensions: 
    - health_check
    - pprof
    - zpages
    ```

    With this configuration the collector will scrape its own metrics.

## Create a new release (with binaries)

Releases are managed by the GH Action workflow. Please commit all changes and then push an annotated tag to the repository.

```shell
git tag -a v<version> -m "<comment>"
git push --tags
```
