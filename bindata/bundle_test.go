/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bindata

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBundle(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	bundle := new(Bundle)
	bundle.PackageName = "bindata"
	err := bundle.Start(buffer)
	assert.Nil(err)
	err = bundle.ProcessPath(buffer, PathConfig{Path: "./testdata/css", Recursive: true})
	assert.Nil(err)
	err = bundle.ProcessPath(buffer, PathConfig{Path: "./testdata/js/app.js", Recursive: false})
	assert.Nil(err)
	err = bundle.Finish(buffer)
	assert.Nil(err)

	assert.NotEmpty(buffer.Bytes())

	assert.Contains(buffer.String(), "package bindata")
	assert.Contains(buffer.String(), "testdata/js/app.js")
	assert.Contains(buffer.String(), "testdata/css/site.css")

	ast, err := parser.ParseFile(token.NewFileSet(), "bindata.go", buffer.Bytes(), parser.ParseComments|parser.AllErrors)
	assert.Nil(err)
	assert.NotNil(ast)
	assert.Len(ast.Imports, 5)
}
