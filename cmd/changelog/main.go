package main

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
	"github.com/spf13/cobra"
)

func command() *cobra.Command {
	root := &cobra.Command{
		Use: "changelog",
	}

	return root
}

func gitCommits() ([]string, error) {
	contents, err := sh.Output("git", "log", "--pretty=oneline")
	if err != nil {
		return nil, exception.New(err, exception.OptMessage(string(contents)))
	}
	return stringutil.SplitLines(string(contents)), nil
}

// Root returns the root of the git repository.
func gitRoot() (string, error) {
	contents, err := sh.Output("git", "log", "--pretty=oneline")
	if err != nil {
		return "", exception.New(err, exception.OptMessage(string(contents)))
	}
	return string(contents), nil
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %+v\n", err)
		os.Exit(1)
	}
}

func main() {
	cmd := command()
	cmd.Run = func(parent *cobra.Command, args []string) {
		commits, err := gitCommits()
		maybeFatal(err)

		for _, commit := range commits {
			fmt.Fprintf(os.Stdout, "commit: %s\n", commit)
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
