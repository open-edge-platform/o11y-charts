// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	cfgType = component.MustNewType("context")
)

var (
	processorCapabilities                  = consumer.Capabilities{MutatesData: true}
	_                     component.Config = (*Config)(nil)
)

// Note: This isn't a valid configuration because the processor would do no work.
func createDefaultConfig() component.Config {
	return &Config{}
}

// NewFactory returns a new factory for the Resource processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		cfgType,
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha),
		processor.WithLogs(createLogsProcessor, component.StabilityLevelAlpha),
		processor.WithTraces(createTracesProcessor, component.StabilityLevelAlpha),
	)
}

type metricsProcessor struct {
	component.Component
	consumer.Metrics
}

func createMetricsProcessor(_ context.Context, set processor.Settings, cfg component.Config, nextConsumer consumer.Metrics) (processor.Metrics, error) {
	ctxtpActions := cfg.(*Config).ActionsConfig
	tracing := trace.WithAttributes(attribute.String("processor", set.ID.String()))
	ctxtp, err := NewContextMetricsProcessor(set.Logger, nextConsumer, tracing, ctxtpActions)
	if err != nil {
		return nil, err
	}
	metricsConsumer, err := consumer.NewMetrics(
		ctxtp.ConsumeMetrics,
		consumer.WithCapabilities(processorCapabilities),
	)
	if err != nil {
		return nil, err
	}
	return &metricsProcessor{
		Component: ctxtp,
		Metrics:   metricsConsumer,
	}, nil
}

type logsProcessor struct {
	component.Component
	consumer.Logs
}

func createLogsProcessor(_ context.Context, set processor.Settings, cfg component.Config, nextConsumer consumer.Logs) (processor.Logs, error) {
	ctxtpActions := cfg.(*Config).ActionsConfig
	tracing := trace.WithAttributes(attribute.String("processor", set.ID.String()))
	ctxtp, err := NewContextLogsProcessor(set.Logger, nextConsumer, tracing, ctxtpActions)
	if err != nil {
		return nil, err
	}
	logsConsumer, err := consumer.NewLogs(
		ctxtp.ConsumeLogs,
		consumer.WithCapabilities(processorCapabilities),
	)
	if err != nil {
		return nil, err
	}
	return &logsProcessor{
		Component: ctxtp,
		Logs:      logsConsumer,
	}, nil
}

type tracesProcessor struct {
	component.Component
	consumer.Traces
}

func createTracesProcessor(_ context.Context, set processor.Settings, cfg component.Config, nextConsumer consumer.Traces) (processor.Traces, error) {
	ctxtpActions := cfg.(*Config).ActionsConfig
	tracing := trace.WithAttributes(attribute.String("processor", set.ID.String()))
	ctxtp, err := NewContextTracesProcessor(set.Logger, nextConsumer, tracing, ctxtpActions)
	if err != nil {
		return nil, err
	}
	tracesConsumer, err := consumer.NewTraces(
		ctxtp.ConsumeTraces,
		consumer.WithCapabilities(processorCapabilities),
	)
	if err != nil {
		return nil, err
	}
	return &tracesProcessor{
		Component: ctxtp,
		Traces:    tracesConsumer,
	}, nil
}
