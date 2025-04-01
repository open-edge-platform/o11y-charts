// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

package contextprocessor

import (
	"errors"
)

var (
	errMissingActionConfig       = errors.New("missing actions configuration")
	errMissingActionConfigKey    = errors.New("missing action key")
	errMissingActionConfigSource = errors.New("missing action source, must be 'from_attribute' or 'value'")
	errMissingActionDeleteParams = errors.New("action delete does not support 'from_attribute' and/or 'value'")
)

// Config represents the receiver config settings within the collector's config.yaml.
type Config struct {
	ActionsConfig []ActionConfig `mapstructure:"actions"`
}

// ActionValue is the enum to capture the four types of actions to perform on the context.
type ActionType string

const (
	// INSERT inserts the new header if it does not exist.
	INSERT ActionType = "insert"
	// UPDATE updates the header value if it exists.
	UPDATE ActionType = "update"
	// UPSERT inserts a header if it does not exist and updates the header if it exists.
	UPSERT ActionType = "upsert"
	// DELETE deletes the header.
	DELETE ActionType = "delete"
)

type ActionConfig struct {
	Key           *string    `mapstructure:"key"`
	Action        ActionType `mapstructure:"action"`
	ValueDefault  *string    `mapstructure:"value"`
	FromAttribute *string    `mapstructure:"from_attribute"`
}

// Validate checks if the extension configuration is valid.
func (cfg *Config) Validate() error {
	if len(cfg.ActionsConfig) == 0 {
		return errMissingActionConfig
	}
	for _, action := range cfg.ActionsConfig {
		if action.Key == nil || *action.Key == "" {
			return errMissingActionConfigKey
		}
		if action.Action != DELETE {
			if action.FromAttribute == nil && action.ValueDefault == nil {
				return errMissingActionConfigSource
			}
		} else {
			if action.FromAttribute != nil || action.ValueDefault != nil {
				return errMissingActionDeleteParams
			}
		}
	}
	return nil
}
