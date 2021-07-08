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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/blend/go-sdk/diff"
	"github.com/blend/go-sdk/stringutil"
)

// New creates a new copyright engine with a given set of config options.
func New(options ...Option) *Copyright {
	var c Copyright
	for _, option := range options {
		option(&c)
	}
	return &c
}

// Copyright is the main type that injects, removes and verifies copyright headers.
type Copyright struct {
	Config // Config holds the configuration opitons.

	// Stdout is the writer for Verbose and Debug output.
	// If it is unset, `os.Stdout` will be used.
	Stdout io.Writer
	// Stderr is the writer for Error output.
	// If it is unset, `os.Stderr` will be used.
	Stderr io.Writer
}

// Inject inserts the copyright header in any matching files that don't already
// have the copyright header.
func (c Copyright) Inject(ctx context.Context) error {
	return c.Walk(ctx, c.inject)
}

// Remove removes the copyright header in any matching files that
// have the copyright header.
func (c Copyright) Remove(ctx context.Context) error {
	return c.Walk(ctx, c.remove)
}

// Verify asserts that the files found during walk
// have the copyright header.
func (c Copyright) Verify(ctx context.Context) error {
	return c.Walk(ctx, c.verify)
}

// Walk traverses the tree recursively from "." and applies the given action.
func (c Copyright) Walk(ctx context.Context, action Action) error {
	noticeBody, err := c.compileNoticeBodyTemplate(c.NoticeBodyTemplateOrDefault())
	if err != nil {
		return err
	}

	c.Verbosef("using excludes: %s", strings.Join(c.Config.Excludes, ", "))
	c.Verbosef("using include files: %s", strings.Join(c.Config.IncludeFiles, ", "))
	c.Verbosef("using notice body:\n%s", noticeBody)

	var didFail bool
	err = filepath.Walk(c.RootOrDefault(), func(path string, info os.FileInfo, fileErr error) error {
		if fileErr != nil {
			return fileErr
		}

		if skipErr := c.includeOrExclude(path, info); skipErr != nil {
			if skipErr == ErrWalkSkip {
				return nil
			}
			return skipErr
		}

		fileExtension := filepath.Ext(path)
		noticeTemplate, ok := c.noticeTemplateByExtension(fileExtension)
		if !ok {
			return fmt.Errorf("invalid copyright injection file; %s", fileExtension)
		}
		notice, err := c.compileNoticeTemplate(noticeTemplate, noticeBody)
		if err != nil {
			return err
		}

		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = action(path, info, fileContents, []byte(notice))
		if err != nil {
			if err == ErrFailure {
				didFail = true
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		if err == ErrFatal {
			// swap out to failure
			return ErrFailure
		}
		return err
	}
	if didFail {
		return ErrFailure
	}
	return nil
}

func (c Copyright) includeOrExclude(path string, info os.FileInfo) error {
	if info.IsDir() {
		if path == c.RootOrDefault() {
			return ErrWalkSkip
		}
	}

	if c.Config.Excludes != nil {
		for _, exclude := range c.Config.Excludes {
			if stringutil.Glob(path, exclude) {
				c.Debugf("path %s matches exclude %s", path, exclude)
				if info.IsDir() {
					return filepath.SkipDir
				}
				return ErrWalkSkip
			}
		}
	}

	if c.Config.IncludeFiles != nil {
		var includePath bool
		for _, include := range c.Config.IncludeFiles {
			if stringutil.Glob(path, include) {
				includePath = true
				break
			}
		}
		if !includePath {
			c.Debugf("path %s does not match any includes", path)
			return ErrWalkSkip
		}
	}
	return nil
}

// noticeTemplateByExtension gets a notice template by extension or the default.
func (c Copyright) noticeTemplateByExtension(fileExtension string) (noticeTemplate string, ok bool) {
	if !strings.HasPrefix(fileExtension, ".") {
		fileExtension = "." + fileExtension
	}

	// check if there is a filetype specific notice template
	extensionNoticeTemplates := c.ExtensionNoticeTemplatesOrDefault()
	if noticeTemplate, ok = extensionNoticeTemplates[fileExtension]; ok {
		return
	}

	// check if we have a default notice template
	if c.NoticeTemplate != "" {
		noticeTemplate = c.NoticeTemplate
		ok = true
		return
	}

	// fail
	return
}

// GetStdout returns standard out.
func (c Copyright) GetStdout() io.Writer {
	if c.QuietOrDefault() {
		return ioutil.Discard
	}
	if c.Stdout != nil {
		return c.Stdout
	}
	return os.Stdout
}

// GetStderr returns standard error.
func (c Copyright) GetStderr() io.Writer {
	if c.QuietOrDefault() {
		return ioutil.Discard
	}
	if c.Stderr != nil {
		return c.Stderr
	}
	return os.Stderr
}

// Verbosef writes to stdout if the `Verbose` flag is true.
func (c Copyright) Verbosef(format string, args ...interface{}) {
	if !c.VerboseOrDefault() {
		return
	}
	fmt.Fprintf(c.GetStdout(), format+"\n", args...)
}

// Debugf writes to stdout if the `Debug` flag is true.
func (c Copyright) Debugf(format string, args ...interface{}) {
	if !c.DebugOrDefault() {
		return
	}
	fmt.Fprintf(c.GetStdout(), format+"\n", args...)
}

//
// actions
//

func (c Copyright) inject(path string, info os.FileInfo, file, notice []byte) error {
	injectedContents := c.injectedContents(path, file, notice)
	if injectedContents == nil {
		return nil
	}
	return ioutil.WriteFile(path, injectedContents, info.Mode().Perm())
}

func (c Copyright) remove(path string, info os.FileInfo, file, notice []byte) error {
	removedContents := c.removedContents(path, file, notice)
	if removedContents == nil {
		return nil
	}
	return ioutil.WriteFile(path, removedContents, info.Mode().Perm())
}

func (c Copyright) verify(path string, _ os.FileInfo, file, notice []byte) error {
	fileExtension := filepath.Ext(path)
	var err error
	if fileExtension == ExtensionGo { // we have to treat go files specially because of build tags
		err = c.goVerifyNotice(path, file, notice)
	} else {
		err = c.verifyNotice(path, file, notice)
	}

	if err != nil {
		fmt.Fprintf(c.GetStderr(), "%+v\n", err)
		if c.Config.ShowDiffOrDefault() {
			c.showDiff(path, file, notice)
		}
		if c.Config.ExitFirstOrDefault() {
			return ErrFatal
		}
		return ErrFailure
	}
	return nil
}

//
// internal helpers
//

func (c Copyright) injectedContents(path string, file, notice []byte) []byte {
	fileExtension := filepath.Ext(path)
	var injectedContents []byte
	if fileExtension == ExtensionGo { // we have to treat go files specially because of build tags
		injectedContents = c.goInjectNotice(path, file, notice)
	} else {
		injectedContents = c.injectNotice(path, file, notice)
	}
	return injectedContents
}

func (c Copyright) removedContents(path string, file, notice []byte) []byte {
	fileExtension := filepath.Ext(path)
	var removedContents []byte
	if fileExtension == ExtensionGo { // we have to treat go files specially because of build tags
		removedContents = c.goRemoveNotice(path, file, notice)
	} else {
		removedContents = c.removeNotice(path, file, notice)
	}
	return removedContents
}

// goInjectNotice handles go files differently because they may contain build tags.
func (c Copyright) goInjectNotice(path string, file, notice []byte) []byte {
	goBuildTag := goBuildTagMatch.FindString(string(file))
	file = goBuildTagMatch.ReplaceAll(file, nil)
	if c.fileHasCopyrightHeader(file, notice) {
		return nil
	}

	c.Verbosef("injecting notice: %s", path)
	file = c.removeCopyrightHeader(file, notice)
	return c.mergeFileSections([]byte(goBuildTag), notice, file)
}

func (c Copyright) injectNotice(path string, file, notice []byte) []byte {
	if c.fileHasCopyrightHeader(file, notice) {
		return nil
	}
	c.Verbosef("injecting notice: %s", path)

	// remove any existing notice-ish looking text ...
	file = c.removeCopyrightHeader(file, notice)
	return c.mergeFileSections(notice, file)
}

func (c Copyright) goRemoveNotice(path string, file, notice []byte) []byte {
	goBuildTag := goBuildTagMatch.FindString(string(file))
	file = goBuildTagMatch.ReplaceAll(file, []byte(""))
	if !c.fileHasCopyrightHeader(file, notice) {
		return nil
	}
	c.Verbosef("removing notice: %s", path)
	return c.mergeFileSections([]byte(goBuildTag), c.removeCopyrightHeader(file, notice))
}

func (c Copyright) removeNotice(path string, file, notice []byte) []byte {
	if !c.fileHasCopyrightHeader(file, notice) {
		return nil
	}
	c.Verbosef("removing notice: %s", path)
	return c.removeCopyrightHeader(file, notice)
}

func (c Copyright) goVerifyNotice(path string, file, notice []byte) error {
	c.Debugf("verifying (go): %s", path)
	fileLessTags := goBuildTagMatch.ReplaceAll(file, nil)
	if !c.fileHasCopyrightHeader(fileLessTags, notice) {
		return fmt.Errorf(VerifyErrorFormat, path)
	}
	return nil
}

func (c Copyright) verifyNotice(path string, file, notice []byte) error {
	c.Debugf("verifying: %s", path)
	if !c.fileHasCopyrightHeader(file, notice) {
		return fmt.Errorf(VerifyErrorFormat, path)
	}
	return nil
}

func (c Copyright) createNoticeMatchExpression(notice []byte, trailingSpaceStrict bool) *regexp.Regexp {
	noticeString := string(notice)
	noticeExpr := yearMatch.ReplaceAllString(regexp.QuoteMeta(noticeString), yearExpr)
	noticeExpr = `^(\s*)` + noticeExpr
	if !trailingSpaceStrict {
		// remove trailing space
		noticeExpr = strings.TrimRightFunc(noticeExpr, unicode.IsSpace)
		// match trailing space
		noticeExpr = noticeExpr + `(\s*)`
	}
	return regexp.MustCompile(noticeExpr)
}

func (c Copyright) fileHasCopyrightHeader(fileContents, notice []byte) bool {
	return c.createNoticeMatchExpression(notice, true).Match(fileContents)
}

func (c Copyright) removeCopyrightHeader(fileContents []byte, notice []byte) []byte {
	return c.createNoticeMatchExpression(notice, false).ReplaceAll(fileContents, nil)
}

func (c Copyright) mergeFileSections(sections ...[]byte) []byte {
	var fullLength int
	for _, section := range sections {
		fullLength += len(section)
	}

	combined := make([]byte, fullLength)

	var written int
	for _, section := range sections {
		copy(combined[written:], section)
		written += len(section)
	}
	return combined
}

func (c Copyright) prefix(prefix string, s string) string {
	lines := strings.Split(s, "\n")
	var output []string
	for _, l := range lines {
		output = append(output, prefix+l)
	}
	return strings.Join(output, "\n")
}

func (c Copyright) compileNoticeTemplate(noticeTemplate, noticeBody string) (string, error) {
	return c.processTemplate(noticeTemplate, c.templateViewModel(map[string]interface{}{
		"Notice": noticeBody,
	}))
}

func (c Copyright) templateViewModel(extra ...map[string]interface{}) map[string]interface{} {
	base := map[string]interface{}{
		"Year":    c.YearOrDefault(),
		"Company": c.CompanyOrDefault(),
		"License": c.LicenseOrDefault(),
	}
	for _, m := range extra {
		for key, value := range m {
			base[key] = value
		}
	}
	return base
}

func (c Copyright) compileRestrictionsTemplate(restrictionsTemplate string) (string, error) {
	return c.processTemplate(restrictionsTemplate, c.templateViewModel())
}

func (c Copyright) compileNoticeBodyTemplate(noticeBodyTemplate string) (string, error) {
	restrictions, err := c.compileRestrictionsTemplate(c.RestrictionsOrDefault())
	if err != nil {
		return "", err
	}
	viewModel := c.templateViewModel(map[string]interface{}{
		"Restrictions": restrictions,
	})
	output, err := c.processTemplate(noticeBodyTemplate, viewModel)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c Copyright) processTemplate(text string, viewmodel interface{}) (string, error) {
	tmpl := template.New("output")
	tmpl = tmpl.Funcs(template.FuncMap{
		"prefix": c.prefix,
	})
	compiled, err := tmpl.Parse(text)
	if err != nil {
		return "", err
	}

	output := new(bytes.Buffer)
	if err = compiled.Execute(output, viewmodel); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (c Copyright) showDiff(path string, file, notice []byte) {
	noticeLineCount := len(stringutil.SplitLines(string(notice),
		stringutil.OptSplitLinesIncludeEmptyLines(true),
		stringutil.OptSplitLinesIncludeNewLine(true),
	))
	fileLines := stringutil.SplitLines(string(file),
		stringutil.OptSplitLinesIncludeEmptyLines(true),
		stringutil.OptSplitLinesIncludeNewLine(true),
	)
	if len(fileLines) < noticeLineCount {
		noticeLineCount = len(fileLines)
	}
	fileTruncated := strings.Join(fileLines[:noticeLineCount], "")
	fileDiff := diff.New().Diff(string(notice), fileTruncated, true /*checklines*/)
	prettyDiff := diff.PrettyText(fileDiff)
	if strings.TrimSpace(prettyDiff) != "" {
		fmt.Fprintf(c.GetStderr(), "%s: diff\n", path)
		fmt.Fprintln(c.GetStderr(), prettyDiff)
	}
}
