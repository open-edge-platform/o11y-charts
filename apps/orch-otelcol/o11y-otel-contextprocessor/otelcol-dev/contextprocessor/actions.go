// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file was copied from https://github.com/springernature/o11y-otel-contextprocessor and adapted to meet specific needs.

package contextprocessor

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// The EventContext the current context and the attributes.
type eventContext struct {
	ctx           context.Context
	cliInfo       client.Info
	resourceAttrs pcommon.Map
	newMetadata   map[string][]string
}

// Sets a new context.
func createEventContext(ctx context.Context, attrs pcommon.Map) *eventContext {
	return &eventContext{
		ctx:           ctx,
		cliInfo:       client.FromContext(ctx),
		resourceAttrs: attrs,
		newMetadata:   make(map[string][]string),
	}
}

func (exc *eventContext) getAttrKey(key, def string) string {
	value := def
	v, exists := exc.resourceAttrs.Get(key)
	if exists {
		switch v.Type() {
		case pcommon.ValueTypeStr:
			value = v.Str()
		default:
			value = v.AsString()
		}
	}
	return value
}

func (exc *eventContext) getContextKey(key string) ([]string, bool) {
	if v, exists := exc.newMetadata[key]; exists {
		return v, exists
	}
	value := exc.cliInfo.Metadata.Get(key)
	return value, len(value) != 0
}

func (exc *eventContext) delContextKey(key string) {
	// Warning: when delete a key is only deleted from newMetadata
	// so it is available again from the actual metadata
	delete(exc.newMetadata, key)
}

func (exc *eventContext) setContextKey(key string, value []string) {
	exc.newMetadata[key] = value
}

func (exc *eventContext) getContext() context.Context {
	return client.NewContext(exc.ctx,
		client.Info{
			Metadata: client.NewMetadata(exc.newMetadata),
		})
}

type Action interface {
	execute(*eventContext)
}

func generateAction(action ActionConfig) (Action, error) {
	valueDefault := ""
	if action.ValueDefault != nil {
		valueDefault = *action.ValueDefault
	}
	fromAttribute := ""
	if action.FromAttribute != nil {
		fromAttribute = *action.FromAttribute
	}
	switch action.Action {
	case INSERT:
		return &actionInsert{
			key:      *action.Key,
			value:    valueDefault,
			fromAttr: fromAttribute,
		}, nil
	case UPSERT:
		return &actionUpsert{
			key:      *action.Key,
			value:    valueDefault,
			fromAttr: fromAttribute,
		}, nil
	case UPDATE:
		return &actionUpdate{
			key:      *action.Key,
			value:    valueDefault,
			fromAttr: fromAttribute,
		}, nil
	case DELETE:
		return &actionDelete{
			key: *action.Key,
		}, nil
	default:
		return nil, errors.New("unknown action type")
	}
}

// Concrete actions

type actionInsert struct {
	key      string
	value    string
	fromAttr string
}

func (a *actionInsert) execute(eventContext *eventContext) {
	value := []string{a.value}
	if len(a.fromAttr) > 0 {
		value[0] = eventContext.getAttrKey(a.fromAttr, a.value)
	}
	if currentValue, exists := eventContext.getContextKey(a.key); !exists {
		eventContext.setContextKey(a.key, value)
	} else {
		eventContext.setContextKey(a.key, currentValue)
	}
}

type actionUpsert struct {
	key      string
	value    string
	fromAttr string
}

func (a *actionUpsert) execute(eventContext *eventContext) {
	value := []string{a.value}
	if len(a.fromAttr) > 0 {
		value[0] = eventContext.getAttrKey(a.fromAttr, a.value)
	}
	eventContext.setContextKey(a.key, value)
}

type actionUpdate struct {
	key      string
	value    string
	fromAttr string
}

func (a *actionUpdate) execute(eventContext *eventContext) {
	value := []string{a.value}
	if len(a.fromAttr) > 0 {
		value[0] = eventContext.getAttrKey(a.fromAttr, a.value)
	}
	if v, exists := eventContext.getContextKey(a.key); exists {
		// There are 2 views here, in this one we add the
		// new value to the current list of strings
		eventContext.setContextKey(a.key, append(v, value[0]))
		// Another option is just overwriting the current value
		// eventContext.setContextKey(a.key, value)
	}
}

type actionDelete struct {
	key string
}

func (a *actionDelete) execute(eventContext *eventContext) {
	eventContext.delContextKey(a.key)
}

type ActionsRunner struct {
	actions []Action
}

func NewActionsRunner() *ActionsRunner {
	return &ActionsRunner{
		actions: make([]Action, 0),
	}
}

func (ar *ActionsRunner) AddAction(action ActionConfig) error {
	a, err := generateAction(action)
	if err == nil {
		ar.actions = append(ar.actions, a)
	}
	return err
}

// The executeCommands method executes all the commands one by one.
func (ar *ActionsRunner) Apply(ctx context.Context, attrs pcommon.Map) context.Context {
	eventContext := createEventContext(ctx, attrs)
	for _, a := range ar.actions {
		a.execute(eventContext)
	}
	return eventContext.getContext()
}
