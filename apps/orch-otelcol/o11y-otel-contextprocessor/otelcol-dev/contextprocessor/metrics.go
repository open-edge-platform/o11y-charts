// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextMetricsProcessor struct {
	contextProcessor
	nextConsumer consumer.Metrics
}

func NewContextMetricsProcessor(
	logger *zap.Logger,
	nextConsumer consumer.Metrics,
	eventOptions trace.SpanStartEventOption,
	actions []ActionConfig) (*contextMetricsProcessor, error) {
	aRunner := NewActionsRunner()
	for _, action := range actions {
		if err := aRunner.AddAction(action); err != nil {
			return nil, err
		}
	}
	return &contextMetricsProcessor{
		contextProcessor: contextProcessor{
			logger:        logger,
			actionsRunner: aRunner,
			eventOptions:  eventOptions,
		},
		nextConsumer: nextConsumer,
	}, nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/consumer#Metrics
func (ctxt *contextMetricsProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) (err error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Start processing.", ctxt.eventOptions)
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len() && err == nil; i++ {
		rm := rms.At(i)
		attrs := rm.Resource().Attributes()
		newCtx := ctxt.actionsRunner.Apply(ctx, attrs)
		newMd := pmetric.NewMetrics()
		rm.CopyTo(newMd.ResourceMetrics().AppendEmpty())
		err = ctxt.nextConsumer.ConsumeMetrics(newCtx, newMd)
	}
	span.AddEvent("End processing.", ctxt.eventOptions)
	return err
}
