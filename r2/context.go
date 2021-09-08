/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"context"
	"net/http"
	"strings"
)

type parameterizedPathKey struct{}

// WithParameterizedPath adds a path with named parameters to a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved.
func WithParameterizedPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, parameterizedPathKey{}, path)
}

// GetParameterizedPath gets a path with named parameters off a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved. Relies on OptParameterizedPath being added to a Request.
func GetParameterizedPath(ctx context.Context) string {
	path, _ := ctx.Value(parameterizedPathKey{}).(string)
	return path
}

// GetParameterizedURLString gets a URL string with named route parameters in place of
// the raw path for a request. Useful for outbound request aggregation for
// metrics and tracing when route parameters are involved.
// Relies on the request's context storing the parameterized path, otherwise will default
// to returning the request `URL`'s `String()`.
func GetParameterizedURLString(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	url := req.URL
	path := GetParameterizedPath(req.Context())
	if path == "" {
		return url.String()
	}

	// Stripped down version of "net/url" `URL.String()`
	var buf strings.Builder
	if url.Scheme != "" {
		buf.WriteString(url.Scheme)
		buf.WriteByte(':')
	}
	if url.Scheme != "" || url.Host != "" {
		if url.Host != "" || url.Path != "" {
			buf.WriteString("//")
		}
		if host := url.Host; host != "" {
			buf.WriteString(host)
		}
	}
	if !strings.HasPrefix(path, "/") {
		buf.WriteByte('/')
	}
	buf.WriteString(path)
	return buf.String()
}
