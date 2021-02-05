/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import "time"

// Config holds the runtime configuration option for the copyright engine.
type Config struct {
	// Root is the starting directory for the file walk.
	Root string `yaml:"root"`
	// NoticeBodyTemplate is the notice body template that will be processed and
	// injected to the relevant extension specific notice template.
	NoticeBodyTemplate string `yaml:"noticeBodyTemplate"`
	// Year is the year to insert into the templates as `{{ .Year }}`
	Year int `yaml:"year"`
	// Company is the company name to insert into the templates as `{{ .Company }}`
	Company string `yaml:"company"`
	// License is the open source license to insert into in templates as `{{ .License }}`
	License string `yaml:"openSourceLicense"`

	// Restrictions an optional template to clarify copyright restrictions or
	// visibility modifiers, which is available in the `NoticeBodyTemplate` as `{{ .Restrictions }}`
	Restrictions string `yaml:"restrictionTemplate"`

	// IncludeFiles are a list of file globs to include.
	IncludeFiles []string `yaml:"includeFiles"`
	// ExcludeFiles are a list of file globs to exclude.
	ExcludeFiles []string `yaml:"excludeFiles"`
	// IncludeDirs are a list of directory globs to include.
	IncludeDirs []string `yaml:"includeDirs"`
	// ExcludeDirs are a list of directory globs to exclude.
	ExcludeDirs []string `yaml:"excludeDirs"`

	// NoticeTemplates should be a map between file extension (including dot)
	// to the relevant notice template for the file. It can include a template variable
	// {{ .Notice }} that will insert the compiled `NoticyBodyTemplate`.
	NoticeTemplates map[string]string

	// ExitFirst indicates if we should return after the first failure.
	ExitFirst *bool `yaml:"exitFirst"`
	// Quiet controls whether output is suppressed.
	Quiet *bool `yaml:"quiet"`
	// Verbose controls whether verbose output is shown.
	Verbose *bool `yaml:"verbose"`
	// Debug controls whether debug output is shown.
	Debug *bool `yaml:"debug"`
}

// RootOrDefault returns the walk root or a default.
func (c Config) RootOrDefault() string {
	if c.Root != "" {
		return c.Root
	}
	return "."
}

// NoticeBodyTemplateOrDefault returns the notice body template or a default.
func (c Config) NoticeBodyTemplateOrDefault() string {
	if c.NoticeBodyTemplate != "" {
		return c.NoticeBodyTemplate
	}
	return DefaultNoticeBodyTemplate
}

// YearOrDefault returns the current year or a default.
func (c Config) YearOrDefault() int {
	if c.Year > 0 {
		return c.Year
	}
	return time.Now().UTC().Year()
}

// CompanyOrDefault returns a company name or a default.
func (c Config) CompanyOrDefault() string {
	if c.Company != "" {
		return c.Company
	}
	return DefaultCompany
}

// LicenseOrDefault returns an open source license or a default.
func (c Config) LicenseOrDefault() string {
	if c.License != "" {
		return c.License
	}
	return DefaultOpenSourceLicense
}

// RestrictionsOrDefault returns restrictions or a default.
func (c Config) RestrictionsOrDefault() string {
	if c.Restrictions != "" {
		return c.Restrictions
	}
	return DefaultRestrictionsInternal
}

// IncludeFilesOrDefault returns a glob list or a default.
func (c Config) IncludeFilesOrDefault() []string {
	if c.IncludeFiles != nil {
		return c.IncludeFiles
	}
	return DefaultIncludeFiles
}

// IncludeDirsOrDefault returns a glob list or a default.
func (c Config) IncludeDirsOrDefault() []string {
	if c.IncludeDirs != nil {
		return c.IncludeDirs
	}
	return DefaultIncludeDirs
}

// ExcludeFilesOrDefault returns a glob list or a default.
func (c Config) ExcludeFilesOrDefault() []string {
	if c.ExcludeFiles != nil {
		return c.ExcludeFiles
	}
	return DefaultExcludeFiles
}

// ExcludeDirsOrDefault returns a glob list or a default.
func (c Config) ExcludeDirsOrDefault() []string {
	if c.ExcludeDirs != nil {
		return c.ExcludeDirs
	}
	return DefaultExcludeDirs
}

// NoticeTemplatesOrDefault returns mapping between file extensions (including dot) to
// the notice templates (i.e. how the template should be commented)
func (c Config) NoticeTemplatesOrDefault() map[string]string {
	if c.NoticeTemplates != nil {
		return c.NoticeTemplates
	}
	return DefaultNoticeTemplates
}

// ExitFirstOrDefault returns a value or a default.
func (c Config) ExitFirstOrDefault() bool {
	if c.ExitFirst != nil {
		return *c.ExitFirst
	}
	return false
}

// QuietOrDefault returns a value or a default.
func (c Config) QuietOrDefault() bool {
	if c.Quiet != nil {
		return *c.Quiet
	}
	return false
}

// VerboseOrDefault returns a value or a default.
func (c Config) VerboseOrDefault() bool {
	if c.Verbose != nil {
		return *c.Verbose
	}
	return false
}

// DebugOrDefault returns a value or a default.
func (c Config) DebugOrDefault() bool {
	if c.Debug != nil {
		return *c.Debug
	}
	return false
}
