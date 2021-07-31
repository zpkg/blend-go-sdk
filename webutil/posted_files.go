/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"io/ioutil"
	"net/http"

	"github.com/blend/go-sdk/ex"
)

const (
	// DefaultPostedFilesMaxMemory is the maximum post body size we will typically consume.
	DefaultPostedFilesMaxMemory = 67_108_864 //64mb
)

// PostedFile is a file that has been posted to an hc endpoint.
type PostedFile struct {
	Key      string
	FileName string
	Contents []byte
}

// PostedFilesOptions are options for the PostedFiles function.
type PostedFilesOptions struct {
	MaxMemory          int64
	ParseMultipartForm bool
	ParseForm          bool
}

// PostedFileOption mutates posted file options.
type PostedFileOption func(*PostedFilesOptions)

// OptPostedFilesMaxMemory sets the max memory for the posted files options (defaults to 64mb).
func OptPostedFilesMaxMemory(maxMemory int64) PostedFileOption {
	return func(pfo *PostedFilesOptions) { pfo.MaxMemory = maxMemory }
}

// OptPostedFilesParseMultipartForm sets if we should parse the multipart form for files (defaults to true).
func OptPostedFilesParseMultipartForm(parseMultipartForm bool) PostedFileOption {
	return func(pfo *PostedFilesOptions) { pfo.ParseMultipartForm = parseMultipartForm }
}

// OptPostedFilesParseForm sets if we should parse the post form for files (defaults to false).
func OptPostedFilesParseForm(parseForm bool) PostedFileOption {
	return func(pfo *PostedFilesOptions) { pfo.ParseForm = parseForm }
}

// PostedFiles returns any files posted
func PostedFiles(r *http.Request, opts ...PostedFileOption) ([]PostedFile, error) {
	var files []PostedFile

	options := PostedFilesOptions{
		MaxMemory:          DefaultPostedFilesMaxMemory,
		ParseMultipartForm: true,
		ParseForm:          false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.ParseMultipartForm {
		if err := r.ParseMultipartForm(options.MaxMemory); err != nil {
			return nil, err
		}
		for key := range r.MultipartForm.File {
			fileReader, fileHeader, err := r.FormFile(key)
			if err != nil {
				return nil, ex.New(err)
			}
			bytes, err := ioutil.ReadAll(fileReader)
			if err != nil {
				return nil, ex.New(err)
			}
			files = append(files, PostedFile{Key: key, FileName: fileHeader.Filename, Contents: bytes})
		}
	}
	if options.ParseForm {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		for key := range r.PostForm {
			if fileReader, fileHeader, err := r.FormFile(key); err == nil && fileReader != nil {
				bytes, err := ioutil.ReadAll(fileReader)
				if err != nil {
					return nil, ex.New(err)
				}
				files = append(files, PostedFile{Key: key, FileName: fileHeader.Filename, Contents: bytes})
			}
		}
	}
	return files, nil
}
