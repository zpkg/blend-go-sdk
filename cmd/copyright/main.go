/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/copyright"
)

type flagStrings []string

func (fs flagStrings) String() string {
	return strings.Join(fs, ", ")
}

func (fs *flagStrings) Set(flagValue string) error {
	if flagValue == "" {
		return fmt.Errorf("invalid flag value; is empty")
	}
	*fs = append(*fs, flagValue)
	return nil
}

var (
	flagNoticeTemplate           string
	flagExtensionNoticeTemplates flagStrings
	flagNoticeBodyTemplate       string
	flagCompany                  string
	flagYear                     int
	flagLicense                  string

	flagRestrictions           string
	flagRestrictionsOpenSource bool
	flagRestrictionsInternal   bool

	flagVerify bool
	flagInject bool
	flagRemove bool

	flagIncludeFiles = flagStrings(copyright.DefaultIncludeFiles)
	flagExcludeFiles = flagStrings(copyright.DefaultExcludeFiles)
	flagIncludeDirs  = flagStrings(copyright.DefaultIncludeDirs)
	flagExcludeDirs  = flagStrings(copyright.DefaultExcludeDirs)

	flagExitFirst bool
	flagQuiet     bool
	flagVerbose   bool
	flagDebug     bool
)

func init() {
	flag.BoolVar(&flagQuiet, "quiet", false, "If all output should be suppressed")
	flag.BoolVar(&flagVerbose, "verbose", false, "If verbose output should be shown")
	flag.BoolVar(&flagDebug, "debug", false, "If debug output should be shown")

	flag.BoolVar(&flagExitFirst, "exit-first", false, "If the program should exit on the first verification error")

	flag.StringVar(&flagCompany, "company", "", "The company name to use in templates as {{ .Company }}")
	flag.IntVar(&flagYear, "year", time.Now().UTC().Year(), "The year to use in templates as {{ .Year }}")
	flag.StringVar(&flagLicense, "license", copyright.DefaultOpenSourceLicense, "The license to use in templates as {{ .License }}")

	flag.StringVar(&flagNoticeTemplate, "notice-template", "", "The notice template; will try as a file path first, if not found then as a string literal")
	flag.StringVar(&flagNoticeBodyTemplate, "notice-body-template", copyright.DefaultNoticeBodyTemplate, "The notice body template; will try as a file path first, if not found then as a string literal")
	flag.StringVar(&flagRestrictions, "restrictions", copyright.DefaultRestrictionsInternal, "The restriction template to compile and insert in the notice body template as {{ .Restrictions }}")

	flag.BoolVar(&flagRestrictionsOpenSource, "restrictions-open-source", false, fmt.Sprintf("The restrictions should be the open source defaults (i.e. %q)", copyright.DefaultRestrictionsOpenSource))
	flag.BoolVar(&flagRestrictionsInternal, "restrictions-internal", false, fmt.Sprintf("The restrictions should be the internal defaults (i.e. %q)", copyright.DefaultRestrictionsInternal))

	flag.BoolVar(&flagVerify, "verify", false, "If we should validate notices are present (exclusive with -inject and -remove) (this is the default)")
	flag.BoolVar(&flagInject, "inject", false, "If we should inject the notice (exclusive with -verify and -remove)")
	flag.BoolVar(&flagRemove, "remove", false, "If we should remove the notice (exclusive with -verify and -inject)")

	flag.Var(&flagExtensionNoticeTemplates, "ext", "Extension specific notice template overrides overrides; should be in the form -ext=js=js_template.txt")

	flag.Var(&flagIncludeFiles, "include-file", "Files to include via glob match")
	flag.Var(&flagExcludeFiles, "exclude-file", "Files to exclude via glob match")
	flag.Var(&flagIncludeDirs, "include-dir", "Directories to include via glob match")
	flag.Var(&flagExcludeDirs, "exclude-dir", "Directories to exclude via glob match")

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), `blend source code copyright management cli

Verify, inject or remove copyright notices from files in a given tree.

By default, this tool verifies that copyright notices are present with no flags provided.

To inject headers:

	> copyright --inject

To remove headers:

	> copyright --remove

`,
		)
		oldUsage()
	}

	flag.Parse()
}

func main() {
	ctx := context.Background()

	if flagNoticeBodyTemplate == "" {
		fmt.Fprintln(os.Stderr, "--notice provided is an empty string; cannot continue")
		os.Exit(1)
	}

	var roots []string
	if args := flag.Args(); len(args) > 0 {
		roots = args
	} else {
		roots = []string{"."}
	}

	var restrictions string
	if flagRestrictionsOpenSource {
		restrictions = copyright.DefaultRestrictionsOpenSource
	} else if flagRestrictionsInternal {
		restrictions = copyright.DefaultRestrictionsInternal
	} else {
		restrictions = flagRestrictions
	}

	extensionNoticeTemplates := copyright.DefaultExtensionNoticeTemplates
	for _, extValue := range flagExtensionNoticeTemplates {
		ext, noticeTemplate, err := parseExtensionNoticeBodyTemplate(extValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		extensionNoticeTemplates[ext] = noticeTemplate
	}

	engine := copyright.Copyright{
		Config: copyright.Config{
			NoticeTemplate:           tryReadFile(flagNoticeTemplate),
			NoticeBodyTemplate:       tryReadFile(flagNoticeBodyTemplate),
			Company:                  flagCompany,
			Restrictions:             restrictions,
			Year:                     flagYear,
			License:                  flagLicense,
			ExtensionNoticeTemplates: extensionNoticeTemplates,
			IncludeFiles:             flagIncludeFiles,
			ExcludeFiles:             flagExcludeFiles,
			IncludeDirs:              flagIncludeDirs,
			ExcludeDirs:              flagExcludeDirs,
			ExitFirst:                &flagExitFirst,
			Quiet:                    &flagQuiet,
			Verbose:                  &flagVerbose,
			Debug:                    &flagDebug,
		},
	}

	var action func(context.Context) error
	var actionLabel string

	if flagVerify {
		action = engine.Verify
		actionLabel = "verify"
	} else if flagInject {
		action = engine.Inject
		actionLabel = "inject"
	} else if flagRemove {
		action = engine.Remove
		actionLabel = "remove"
	} else {
		action = engine.Verify
		actionLabel = "verify"
	}

	var didFail bool
	for _, root := range roots {
		engine.Root = root
		maybeFail(ctx, action, &didFail)
	}
	if didFail {
		if !flagQuiet {
			fmt.Printf("copyright %s %s!\n", actionLabel, ansi.Red("failed"))
		}
		os.Exit(1)
	}
	if !flagQuiet {
		fmt.Printf("copyright %s %s!\n", actionLabel, ansi.Green("ok"))
	}
}

func tryReadFile(path string) string {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return path
	}
	return string(contents)
}

func parseExtensionNoticeBodyTemplate(extensionNoticeBodyTemplate string) (extension, noticeBodyTemplate string, err error) {
	parts := strings.SplitN(extensionNoticeBodyTemplate, "=", 2)
	if len(parts) < 2 {
		err = fmt.Errorf("invalid `-ext` value; %s", extensionNoticeBodyTemplate)
		return
	}
	extension = parts[0]
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	noticeBodyTemplate = tryReadFile(parts[1])
	return
}

func maybeFail(ctx context.Context, action func(context.Context) error, didFail *bool) {
	err := action(ctx)
	if err != nil {
		if err == copyright.ErrFailure {
			*didFail = true
			return
		}
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
