package main

import (
	"fmt"
	"os"

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

func getCommits() ([]string, error) {
	contents, err := sh.OutputParsed("git log")
	if err != nil {
		return nil, err
	}
	return stringutil.SplitLines(string(contents)), nil
}

func main() {
	cmd := command()
	cmd.Run = func(parent *cobra.Command, args []string) {
		os.Exit(0)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
