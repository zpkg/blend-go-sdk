/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"errors"
	"regexp"
)

// DefaultCompany is the default company to inject into the notice template.
const DefaultCompany = "Blend Labs, Inc"

// DefaultOpenSourceLicense is the default open source license.
const DefaultOpenSourceLicense = "MIT"

// DefaultRestrictionsInternal are the default copyright restrictions to inject into the notice template.
const DefaultRestrictionsInternal = "Blend Confidential - Restricted"

// DefaultRestrictionsOpenSource are the default open source restrictions.
const DefaultRestrictionsOpenSource = `Use of this source code is governed by a {{ .License }} license that can be found in the LICENSE file.`

// DefaultNoticeBodyTemplate is the default notice body template.
const DefaultNoticeBodyTemplate = `Copyright (c) {{ .Year }} - Present. {{ .Company }}. All rights reserved
{{ .Restrictions }}`

var (
	// DefaultExtensionNoticeTemplates is a mapping between file extension (including the prefix dot) to the notice templates.
	DefaultExtensionNoticeTemplates = map[string]string{
		".css":  cssNoticeTemplate,
		".go":   goNoticeTemplate,
		".html": htmlNoticeTemplate,
		".js":   jsNoticeTemplate,
		".jsx":  jsNoticeTemplate,
		".py":   pythonNoticeTemplate,
		".sass": sassNoticeTemplate,
		".scss": scssNoticeTemplate,
		".ts":   tsNoticeTemplate,
		".tsx":  tsNoticeTemplate,
		".yaml": yamlNoticeTemplate,
		".yml":  yamlNoticeTemplate,
	}

	// DefaultIncludeFiles is the default included files list.
	DefaultIncludeFiles = []string{
		"*.css",
		"*.go",
		"*.html",
		"*.js",
		"*.jsx",
		"*.py",
		"*.sass",
		"*.scss",
		"*.ts",
		"*.tsx",
		"*.yaml",
		"*.yml",
	}

	// DefaultIncludeDirs is the default included directories.
	DefaultIncludeDirs = []string{
		"*",
	}

	// DefaultExcludeFiles is the default excluded files list.
	DefaultExcludeFiles = []string{}

	// DefaultExcludeDirs is the default excluded directories.
	DefaultExcludeDirs = []string{
		".git/*",
		".github/*",
		"*/_config",
		"*/_config/*",
		"*/dist/*",
		"*/node_modules/*",
		"*/testdata",
		"*/testdata/*",
		"*/vendor/*",
		"protogen/*",
		"vendor/*",
		"venv/*",
	}
)

// Errors
var (
	verifyErrorFormat = "%s: file copyright header missing or invalid; please use `copyright --inject` to add it"
	ErrFailure        = errors.New("failure; one or more steps failed")
)

const (
	// goNoticeTemplate is the notice template specific to go files
	// note: it _must_ end in two newlines to prevent linting / compiler failures.
	goNoticeTemplate = `/*

{{ .Notice }}

*/

`

	yamlNoticeTemplate = `#
{{ .Notice | prefix "# " }}
#
`

	htmlNoticeTemplate = `<!--
{{ .Notice }}
-->
`

	jsNoticeTemplate = `/**
{{ .Notice | prefix " * " }}
 */
`

	tsNoticeTemplate = `/**
{{ .Notice | prefix " * " }}
 */
`

	cssNoticeTemplate = `/*
{{ .Notice | prefix " * " }}
 */
`

	scssNoticeTemplate = `/*
{{ .Notice | prefix " * " }}
 */
`

	sassNoticeTemplate = `/*
{{ .Notice | prefix " * " }}
 */
`

	pythonNoticeTemplate = `'''
{{ .Notice }}
'''
`
)

const (
	goBuildTagExpr = `(?s)^\/\/ \+build([^\n]+)(\n{2})`
	yearExpr       = `([0-9]{4,}?)`
)

var (
	goBuildTagMatch = regexp.MustCompile(goBuildTagExpr)
	yearMatch       = regexp.MustCompile(yearExpr)
)
