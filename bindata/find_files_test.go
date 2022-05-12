/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bindata

import (
	"regexp"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBundleFindFiles(t *testing.T) {
	assert := assert.New(t)

	var files []*File
	err := (&Bundle{}).findFiles("./testdata", nil, func(f *File) error {
		files = append(files, f)
		return nil
	})
	assert.Nil(err)
	assert.Len(files, 2)
	assert.Equal("testdata/css/site.css", files[0].Name)
	assert.Equal("testdata/js/app.js", files[1].Name)
}

func TestBundleFindFilesIgnores(t *testing.T) {
	assert := assert.New(t)

	ignoreCSS := regexp.MustCompile(".css$")

	var files []*File
	err := (&Bundle{}).findFiles("./testdata", []*regexp.Regexp{ignoreCSS}, func(f *File) error {
		files = append(files, f)
		return nil
	})
	assert.Nil(err)
	assert.Len(files, 1)
	assert.Equal("testdata/js/app.js", files[0].Name)
}
