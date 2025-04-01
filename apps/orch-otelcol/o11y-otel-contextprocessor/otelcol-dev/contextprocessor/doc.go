// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

//go:generate mdatagen metadata.yaml

// Package contextprocessor contains the logic to manage context metadata
// It supports insert, update, upsert and delete as actions.
package contextprocessor // import "github.com/jriguera/opentelemetry-collector-contrib/processor/contextprocessor"
