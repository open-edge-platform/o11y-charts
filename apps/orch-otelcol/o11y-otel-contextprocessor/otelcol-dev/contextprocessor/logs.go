// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextLogsProcessor struct {
	contextProcessor
	nextConsumer consumer.Logs
}

func NewContextLogsProcessor(
	logger *zap.Logger,
	nextConsumer consumer.Logs,
	eventOptions trace.SpanStartEventOption,
	actions []ActionConfig) (*contextLogsProcessor, error) {
	aRunner := NewActionsRunner()
	for _, action := range actions {
		if err := aRunner.AddAction(action); err != nil {
			return nil, err
		}
	}
	return &contextLogsProcessor{
		contextProcessor: contextProcessor{
			logger:        logger,
			actionsRunner: aRunner,
			eventOptions:  eventOptions,
		},
		nextConsumer: nextConsumer,
	}, nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/consumer#Logs
func (ctxt *contextLogsProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) (err error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Start processing.", ctxt.eventOptions)
	rsl := ld.ResourceLogs()
	for i := 0; i < rsl.Len() && err == nil; i++ {
		rl := rsl.At(i)
		attrs := rl.Resource().Attributes()
		newCtx := ctxt.actionsRunner.Apply(ctx, attrs)
		newLd := plog.NewLogs()
		rl.CopyTo(newLd.ResourceLogs().AppendEmpty())
		err = ctxt.nextConsumer.ConsumeLogs(newCtx, newLd)
	}
	span.AddEvent("End processing.", ctxt.eventOptions)
	return err
}
