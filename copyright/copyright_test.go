/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"
)

func Test_Copyright_GetStdout(t *testing.T) {
	its := assert.New(t)

	c := New()

	its.Equal(os.Stdout, c.GetStdout())
	buf := new(bytes.Buffer)
	c.Stdout = buf
	its.Equal(c.Stdout, c.GetStdout())
	c.Quiet = ref.Bool(true)
	its.Equal(io.Discard, c.GetStdout())
}

func Test_Copyright_GetStderr(t *testing.T) {
	its := assert.New(t)

	c := New()

	its.Equal(os.Stderr, c.GetStderr())
	buf := new(bytes.Buffer)
	c.Stderr = buf
	its.Equal(buf, c.GetStderr())
	c.Quiet = ref.Bool(true)
	its.Equal(io.Discard, c.GetStderr())
}

func Test_Copyright_mergeFileSections(t *testing.T) {
	its := assert.New(t)

	merged := Copyright{}.mergeFileSections([]byte("foo"), []byte("bar"), []byte("baz"))
	its.Equal("foobarbaz", string(merged))
}

func Test_Copyright_fileHasCopyrightHeader(t *testing.T) {
	its := assert.New(t)

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	goodCorpusWithNotice := Copyright{}.mergeFileSections([]byte(notice), goodCorpus)
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2021")
	its.True((Copyright{}).fileHasCopyrightHeader(goodCorpusWithNotice, []byte(notice)))
}

func Test_Copyright_fileHasCopyrightHeader_invalid(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var invalidCorpus = []byte(`foo
bar
baz
`)
	expectedNotice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	its.False(c.fileHasCopyrightHeader(invalidCorpus, []byte(expectedNotice)), "we haven't added the notice")
}

func Test_Copyright_fileHasCopyrightHeader_differentYear(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2020))
	its.Nil(err)

	goodCorpusWithNotice := c.mergeFileSections(notice, goodCorpus)
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2020")

	newNotice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	its.True(c.fileHasCopyrightHeader(goodCorpusWithNotice, []byte(newNotice)))
}

func Test_Copyright_fileHasCopyrightHeader_leadingWhitespace(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	goodCorpusWithNotice := c.mergeFileSections([]byte("\n\n"), notice, goodCorpus)
	its.HasPrefix(string(goodCorpusWithNotice), "\n\n")
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2021")

	its.True(c.fileHasCopyrightHeader(goodCorpusWithNotice, []byte(notice)))
}

func Test_Copyright_goBuildTagMatch(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	buildTag := []byte(`// +build foo

`)
	corpus := []byte(`foo
bar
baz
`)

	file := (Copyright{}).mergeFileSections(buildTag, corpus)

	its.False(goBuildTagMatch.Match(corpus))
	its.True(goBuildTagMatch.Match(c.mergeFileSections(buildTag)))

	found := goBuildTagMatch.FindAll(file, -1)
	its.NotEmpty(found)
	its.True(goBuildTagMatch.Match(file))
}

func Test_Copyright_goInjectNotice(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasSuffix(string(output), string(file))
}

func Test_Copyright_goInjectNotice_buildTags(t *testing.T) {
	its := assert.New(t)
	c := Copyright{}

	buildTag := []byte(`// +build foo`)
	corpus := []byte(`foo
bar
baz
`)

	file := c.mergeFileSections(buildTag, []byte("\n\n"), corpus)

	notice, err := generateGoNotice(OptYear(2021))
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasPrefix(string(output), string(buildTag)+"\n")
	its.HasSuffix(string(output), string(corpus))

	outputRepeat := c.goInjectNotice("foo.go", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func Test_Copyright_injectNotice_typescript(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(`foo
bar
baz
`)

	notice, err := generateTypescriptNotice(OptYear(2021))
	its.Nil(err)

	output := c.injectNotice("foo.ts", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasSuffix(string(output), string(file))

	outputRepeat := c.injectNotice("foo.ts", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func Test_Copyright_goInjectNotice_openSource(t *testing.T) {
	its := assert.New(t)

	c := new(Copyright)

	file := []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(
		OptYear(2021),
		OptLicense("Apache 2.0"),
		OptRestrictions(DefaultRestrictionsOpenSource),
	)
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.Contains(string(output), "Use of this source code is governed by a Apache 2.0")
	its.HasSuffix(string(output), string(file))
}

func generateGoNotice(opts ...Option) ([]byte, error) {
	c := New(opts...)
	noticeBody, err := c.compileNoticeBodyTemplate(c.NoticeBodyTemplateOrDefault())
	if err != nil {
		return nil, err
	}

	compiled, err := c.compileNoticeTemplate(goNoticeTemplate, noticeBody)
	if err != nil {
		return nil, err
	}
	return []byte(compiled), nil
}

func generateTypescriptNotice(opts ...Option) ([]byte, error) {
	c := New(opts...)
	noticeBody, err := c.compileNoticeBodyTemplate(c.NoticeBodyTemplateOrDefault())
	if err != nil {
		return nil, err
	}

	compiled, err := c.compileNoticeTemplate(tsNoticeTemplate, noticeBody)
	if err != nil {
		return nil, err
	}
	return []byte(compiled), nil
}

func Test_Copyright_GetNoticeTemplate(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	noticeTemplate, ok := c.noticeTemplateByExtension(".js")
	its.True(ok)
	its.Equal(jsNoticeTemplate, noticeTemplate)

	// it handles no dot prefix
	noticeTemplate, ok = c.noticeTemplateByExtension("js")
	its.True(ok)
	its.Equal(jsNoticeTemplate, noticeTemplate)

	// it handles another file type
	noticeTemplate, ok = c.noticeTemplateByExtension(".go")
	its.True(ok)
	its.Equal(goNoticeTemplate, noticeTemplate)

	noticeTemplate, ok = c.noticeTemplateByExtension("not-a-real-extension")
	its.False(ok)
	its.Empty(noticeTemplate)

	withDefault := Copyright{
		Config: Config{
			NoticeTemplate: "this is just a test",
		},
	}

	noticeTemplate, ok = withDefault.noticeTemplateByExtension("not-a-real-extension")
	its.True(ok)
	its.Equal("this is just a test", noticeTemplate)
}

type mockInfoDir string

func (mid mockInfoDir) Name() string       { return string(mid) }
func (mid mockInfoDir) Size() int64        { return 1 << 8 }
func (mid mockInfoDir) Mode() fs.FileMode  { return fs.FileMode(0755) }
func (mid mockInfoDir) ModTime() time.Time { return time.Now().UTC() }
func (mid mockInfoDir) IsDir() bool        { return true }
func (mid mockInfoDir) Sys() interface{}   { return nil }

type mockInfoFile string

func (mif mockInfoFile) Name() string       { return string(mif) }
func (mif mockInfoFile) Size() int64        { return 1 << 8 }
func (mif mockInfoFile) Mode() fs.FileMode  { return fs.FileMode(0755) }
func (mif mockInfoFile) ModTime() time.Time { return time.Now().UTC() }
func (mif mockInfoFile) IsDir() bool        { return false }
func (mif mockInfoFile) Sys() interface{}   { return nil }

func Test_Copyright_includeOrExclude(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	testCases := [...]struct {
		Config   Config
		Path     string
		Info     fs.FileInfo
		Expected error
	}{
		/*0*/ {Config: Config{Root: "."}, Path: ".", Info: mockInfoDir("."), Expected: ErrWalkSkip},
		/*1*/ {Config: Config{Root: ".", Excludes: []string{"/foo/**"}}, Path: "/foo/bar", Info: mockInfoDir("bar"), Expected: filepath.SkipDir},
		/*2*/ {Config: Config{Root: ".", Excludes: []string{"/foo/**"}}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: ErrWalkSkip},
		/*3*/ {Config: Config{Root: ".", IncludeFiles: []string{"/foo/bar/*.jpg"}}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: nil},
		/*4*/ {Config: Config{Root: ".", Excludes: []string{}, IncludeFiles: []string{}}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: ErrWalkSkip},
		/*5*/ {Config: Config{Root: "."}, Path: "/foo/bar/baz.jpg", Info: mockInfoFile("baz.jpg"), Expected: nil},
		/*6*/ {Config: Config{Root: "."}, Path: "/foo/bar/baz.jpg", Info: mockInfoDir("baz"), Expected: ErrWalkSkip},
	}

	for index, tc := range testCases {
		c := Copyright{Config: tc.Config}
		its.Equal(tc.Expected, c.includeOrExclude(tc.Path, tc.Info), fmt.Sprintf("test %d", index))
	}
}

const (
	tsFile0 = `import * as axios from 'axios';`
	pyFile0 = `from __future__ import print_function
	
		import logging
		import os
		import shutil
		import sys
		import requests
		import uuid
		import json`

	goFile0 = `// +build tools
		package tools
		
		import (
			// goimports organizes imports for us
			_ "golang.org/x/tools/cmd/goimports"
		
			// golint is an opinionated linter
			_ "golang.org/x/lint/golint"
		
			// ineffassign is an opinionated linter
			_ "github.com/gordonklaus/ineffassign"
		
			// staticcheck is ineffassign but better
			_ "honnef.co/go/tools/cmd/staticcheck"
		)
		`
)

func createTestFS(its *assert.Assertions) (tempDir string, revert func()) {
	// create a temp dir
	var err error
	tempDir, err = ioutil.TempDir("", "coverage_test")
	its.Nil(err)
	revert = func() {
		os.RemoveAll(tempDir)
	}

	// create some files
	err = os.MkdirAll(filepath.Join(tempDir, "foo", "bar"), 0755)
	its.Nil(err)
	err = os.MkdirAll(filepath.Join(tempDir, "bar", "foo"), 0755)
	its.Nil(err)

	err = os.MkdirAll(filepath.Join(tempDir, "not-bar", "not-foo"), 0755)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "foo", "bar", "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "foo", "bar", "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "foo", "bar", "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "bar", "foo", "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "bar", "foo", "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "bar", "foo", "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file0.py"), []byte(pyFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file0.ts"), []byte(tsFile0), 0644)
	its.Nil(err)

	err = ioutil.WriteFile(filepath.Join(tempDir, "not-bar", "not-foo", "file0.go"), []byte(goFile0), 0644)
	its.Nil(err)
	return
}

func Test_Copyright_Walk(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptRoot(tempDir),
		OptIncludeFiles("*.py", "*.ts"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)

	var err error
	var seen []string
	err = c.Walk(context.TODO(), func(path string, info os.FileInfo, file, notice []byte) error {
		seen = append(seen, path)
		return nil
	})
	its.Nil(err)
	expected := []string{
		filepath.Join(tempDir, "bar", "foo", "file0.py"),
		filepath.Join(tempDir, "bar", "foo", "file0.ts"),
		filepath.Join(tempDir, "file0.py"),
		filepath.Join(tempDir, "file0.ts"),
		filepath.Join(tempDir, "foo", "bar", "file0.py"),
		filepath.Join(tempDir, "foo", "bar", "file0.ts"),
	}
	its.Equal(expected, seen)
}

func Test_Copyright_Inject(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptRoot(tempDir),
		OptIncludeFiles("*.py", "*.ts", "*.go"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)

	err := c.Inject(context.TODO())
	its.Nil(err)

	contents, err := ioutil.ReadFile(filepath.Join(tempDir, "bar", "foo", "file0.py"))
	its.Nil(err)
	its.HasPrefix(string(contents), "#\n# Copyright")

	contents, err = ioutil.ReadFile(filepath.Join(tempDir, "bar", "foo", "file0.ts"))
	its.Nil(err)
	its.HasPrefix(string(contents), "/**\n * Copyright")
}

func Test_Copyright_Verify(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptRoot(tempDir),
		OptIncludeFiles("*.py", "*.ts", "*.go"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)
	c.Stdout = new(bytes.Buffer)
	c.Stderr = new(bytes.Buffer)

	err := c.Verify(context.TODO())
	its.NotNil(err)

	err = c.Inject(context.TODO())
	its.Nil(err)

	err = c.Verify(context.TODO())
	its.Nil(err)
}

func Test_Copyright_Remove(t *testing.T) {
	its := assert.New(t)

	tempDir, revert := createTestFS(its)
	defer revert()

	c := New(
		OptRoot(tempDir),
		OptIncludeFiles("*.py", "*.ts", "*.go"),
		OptExcludes("*/not-bar/*", "*/not-foo/*"),
	)
	c.Stdout = new(bytes.Buffer)
	c.Stderr = new(bytes.Buffer)

	err := c.Inject(context.TODO())
	its.Nil(err)

	err = c.Verify(context.TODO())
	its.Nil(err)

	err = c.Remove(context.TODO())
	its.Nil(err)

	err = c.Verify(context.TODO())
	its.NotNil(err)
}
