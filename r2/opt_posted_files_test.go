/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestOptPostedFiles(t *testing.T) {
	its := assert.New(t)

	r := New(TestURL, OptPostedFiles(
		webutil.PostedFile{
			Key:      "form-key",
			FileName: "file.txt",
			Contents: []byte("this is a test"),
		},
		webutil.PostedFile{
			Key:      "form-key-2",
			FileName: "file2.txt",
			Contents: []byte("this is a test2"),
		},
	))
	its.NotNil(r.Request.Body)

	files, err := webutil.PostedFiles(r.Request)
	its.Nil(err)
	its.Len(files, 2)

	its.Equal("form-key", files[0].Key)
	its.Equal("file.txt", files[0].FileName)
	its.Equal("this is a test", string(files[0].Contents))

	its.Equal("form-key-2", files[1].Key)
	its.Equal("file2.txt", files[1].FileName)
	its.Equal("this is a test2", string(files[1].Contents))
}
