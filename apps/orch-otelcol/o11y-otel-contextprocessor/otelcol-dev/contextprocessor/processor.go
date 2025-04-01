// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextProcessor struct {
	logger        *zap.Logger
	actionsRunner *ActionsRunner
	cancel        context.CancelFunc
	eventOptions  trace.SpanStartEventOption
}

// Start implements https://pkg.go.dev/go.opentelemetry.io/collector/component#Component.
func (ctxt *contextProcessor) Start(_ context.Context, host component.Host) error {
	_, ctxt.cancel = context.WithCancel(context.Background())
	for k := range host.GetExtensions() {
		ctxt.logger.Info("Extension", zap.String("id", k.String()))
	}
	return nil
}

// Shutdown implements https://pkg.go.dev/go.opentelemetry.io/collector/component#Component.
func (ctxt *contextProcessor) Shutdown(context.Context) error {
	ctxt.cancel()
	return nil
}
