/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"context"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sanitize"
)

// NewAddListenerOptions creates a new add listener options.
func NewAddListenerOptions(opts ...AddListenerOption) AddListenerOptions {
	var options AddListenerOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// AddListenerOptions are options for adding listeners.
type AddListenerOptions struct {
	IncludeLoggerLabelsAsTags bool
	RequestSanitizeDefaults   []sanitize.RequestOption
}

// GetLoggerTags gets the logger tags from a context if they're set to be included.
func (options AddListenerOptions) GetLoggerTags(ctx context.Context) (tags []string) {
	if options.IncludeLoggerLabelsAsTags {
		// append the labels
		if labels := logger.GetLabels(ctx); len(labels) > 0 {
			for key, value := range labels {
				tags = append(tags, Tag(key, value))
			}
		}
	}
	return
}

// OptIncludeLoggerLabelsAsTags includes logger labels as tags.
func OptIncludeLoggerLabelsAsTags(include bool) AddListenerOption {
	return func(a *AddListenerOptions) { a.IncludeLoggerLabelsAsTags = include }
}

// OptRequestSanitizeDefaults includes logger labels as tags.
func OptRequestSanitizeDefaults(opts ...sanitize.RequestOption) AddListenerOption {
	return func(a *AddListenerOptions) { a.RequestSanitizeDefaults = opts }
}

// AddListenerOption mutates AddListenerOptions
type AddListenerOption func(*AddListenerOptions)
