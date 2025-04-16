<!--
SPDX-FileCopyrightText: (C) 2025 Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# Edge Orchestrator Observability Charts

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[Edge Node Observability Stack]: charts/edgenode-observability
[Edge Node Observability Dashboards]: charts/edgenode-dashboards
[Orchestrator Observability Stack]: charts/orchestrator-observability
[Orchestrator Observability Dashboards]: charts/orchestrator-dashboards

[Grafana Proxy]: apps/grafana-proxy
[OpenTelemetry Collector]: apps/orch-otelcol

[Documentation]: https://docs.openedgeplatform.intel.com/edge-manage-docs/main/developer_guide/observability/arch/index.html
[Contributor's Guide]: https://docs.openedgeplatform.intel.com/edge-manage-docs/main/developer_guide/contributor_guide/index.html
[Edge Orchestrator Community]: https://github.com/open-edge-platform
[Troubleshooting]: https://docs.openedgeplatform.intel.com/edge-manage-docs/main/developer_guide/troubleshooting/index.html
[Contact us]: https://github.com/open-edge-platform

[Apache 2.0 License]: LICENSES/Apache-2.0.txt

## Overview

Edge Orchestrator Observability Charts provide key Observability charts deployable on Edge Orchestrator:

- [Edge Node Observability Stack] along with [Edge Node Observability Dashboards]
- [Orchestrator Observability Stack] along with [Orchestrator Observability Dashboards]

Additionally, the charts utilize customized applications to handle multi-tenancy aspects:

- [Grafana Proxy]
- [OpenTelemetry Collector]

Read more about Edge Orchestrator Observability Charts in the [Documentation].

## Get Started

This repository consists of several projects located in the `apps` and `charts` directories.
If these commands are run in the main repository directory, they will be executed for all projects.
To run them for a specific project, execute the commands in the project's directory, e.g., `charts/edgenode-dashboards`.

To set up the development environment and work on this project, follow the steps below.
All necessary tools will be installed using the `install-tools` target.
Note that `docker` and `asdf` must be installed beforehand.

## Develop

The code of this project is maintained and released in CI using the `VERSION` file.
In addition, the chart is versioned with the same tag as the `VERSION` file.

This is mandatory to keep all chart versions and app versions coherent.

To bump the version, increment the version in the `VERSION` file and run the following command
(to set `version` and `appVersion` in the `Chart.yaml` automatically):

```sh
make helm-build
```

### Install Tools

To install all the necessary tools needed for development the project, run:

```sh
make install-tools
```

### Build

To build the project, use the following command:

```sh
make build
```

### Lint

To lint the code and ensure it adheres to the coding standards, run:

```sh
make lint
```

### Test

To run the tests and verify the functionality of the project, use:

```sh
make test
```

### Docker Build

To build the Docker image for the project, run:

```sh
make docker-build
```

### Helm Build

To package the Helm chart for the project, use:

```sh
make helm-build
```

### Docker Push

To push the Docker image to the registry, run:

```sh
make docker-push
```

### Helm Push

To push the Helm chart to the repository, use:

```sh
make helm-push
```

### Kind All

To load the Docker image into a local Kind cluster, run:

```sh
make kind-all
```

## Contribute

To learn how to contribute to the project, see the [Contributor's Guide].

## Community and Support

To learn more about the project, its community, and governance, visit the [Edge Orchestrator Community].

For support, start with [Troubleshooting] or [Contact us].

## License

Edge Orchestrator Observability Charts is licensed under [Apache 2.0 License].

Last Updated Date: {March 28, 2025}
