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
	flagNotice  string
	flagCompany string
	flagYear    int

	flagInject bool
	flagRemove bool

	flagIncludeFiles = flagStrings{}
	flagExcludeFiles = flagStrings{}
	flagIncludeDirs  = flagStrings{}
	flagExcludeDirs  = flagStrings{}

	flagExitFirst bool
	flagVerbose   bool
	flagDebug     bool
)

func init() {
	flag.BoolVar(&flagVerbose, "verbose", false, "If verbose output should be shown")
	flag.BoolVar(&flagDebug, "debug", false, "If verbose output should be shown")

	flag.BoolVar(&flagExitFirst, "exit-first", false, "If the program should exit on the first verification error")

	flag.StringVar(&flagNotice, "notice", copyright.DefaultNoticeBodyTemplate, "The notice body template; use `-` to read from standard input")
	flag.StringVar(&flagCompany, "company", copyright.DefaultCompany, "The company name to use in the notice body template")
	flag.IntVar(&flagYear, "year", time.Now().UTC().Year(), "The year to use in the notice body template")

	flag.BoolVar(&flagInject, "inject", false, "If we should inject the notice")
	flag.BoolVar(&flagRemove, "remove", false, "If we should remove the notice")

	flag.Var(&flagIncludeFiles, "include-file", "Files to include via glob match")
	flag.Var(&flagExcludeFiles, "exclude-file", "Files to exclude via glob match")
	flag.Var(&flagIncludeDirs, "include-dir", "Directories to include via glob match")
	flag.Var(&flagExcludeDirs, "exclude-dir", "Directories to exclude via glob match")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	if flagNotice == "" {
		fmt.Fprintln(os.Stderr, "--notice was provided an empty string; cannot continue")
		os.Exit(1)
	}

	if strings.TrimSpace(flagNotice) == "-" {
		notice, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		flagNotice = string(notice)
	}

	var roots []string
	if args := flag.Args(); len(args) > 0 {
		roots = args
	} else {
		roots = []string{"."}
	}

	var didFail bool
	for _, root := range roots {
		engine := copyright.Copyright{
			Config: copyright.Config{
				Root:               root,
				NoticeBodyTemplate: flagNotice,
				Company:            flagCompany,
				Year:               flagYear,
				IncludeFiles:       flagStringsWithDefault(flagIncludeFiles, copyright.DefaultIncludeFiles),
				ExcludeFiles:       flagStringsWithDefault(flagExcludeFiles, copyright.DefaultExcludeFiles),
				IncludeDirs:        flagStringsWithDefault(flagIncludeDirs, copyright.DefaultIncludeDirs),
				ExcludeDirs:        flagStringsWithDefault(flagExcludeDirs, copyright.DefaultExcludeDirs),
				ExitFirst:          &flagExitFirst,
				Verbose:            &flagVerbose,
				Debug:              &flagDebug,
			},
		}

		if flagInject {
			maybeFail(ctx, engine.Inject, &didFail)
		}
		if flagRemove {
			maybeFail(ctx, engine.Remove, &didFail)
		}
		maybeFail(ctx, engine.Verify, &didFail)
	}
	if didFail {
		os.Exit(1)
	}
	fmt.Printf("copyright %s!\n", ansi.Green("ok"))
}

func flagStringsWithDefault(flagPointer flagStrings, defaultValues []string) []string {
	if flagPointer != nil && len(flagPointer) > 0 {
		return flagPointer
	}
	return defaultValues
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
