/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"net/http"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// CachedStaticFile is a memory mapped static file.
type CachedStaticFile struct {
	Path     string
	Size     int
	ETag     string
	ModTime  time.Time
	Contents *bytes.Reader
}

// Render implements Result.
func (csf CachedStaticFile) Render(ctx *Ctx) error {
	if csf.ETag != "" {
		ctx.Response.Header().Set(webutil.HeaderETag, csf.ETag)
	}
	http.ServeContent(ctx.Response, ctx.Request, csf.Path, csf.ModTime, csf.Contents)
	return nil
}
