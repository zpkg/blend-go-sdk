/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"

	"github.com/zpkg/blend-go-sdk/ex"
)

// MustParseSameSite parses a string value for same site and panics on error.
func MustParseSameSite(sameSite string) http.SameSite {
	value, err := ParseSameSite(sameSite)
	if err != nil {
		panic(err)
	}
	return value
}

// ParseSameSite parses a string value for same site.
func ParseSameSite(sameSite string) (http.SameSite, error) {
	switch sameSite {
	case SameSiteStrict:
		return http.SameSiteStrictMode, nil
	case SameSiteLax:
		return http.SameSiteLaxMode, nil
	case SameSiteDefault:
		return http.SameSiteDefaultMode, nil
	default:
		return http.SameSite(-1), ex.New(ErrInvalidSameSite, ex.OptMessagef("value: %s", sameSite))
	}
}
