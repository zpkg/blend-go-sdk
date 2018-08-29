package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blend/go-sdk/semver"
)

func usage() {
	fmt.Fprint(os.Stdout, "version validates and manage versions from a given file\n")
	fmt.Fprint(os.Stdout, "\ncommands:\n")
	fmt.Fprint(os.Stdout, "\tvalidate\t\t\t\tvalidate a given version file\n")
	fmt.Fprint(os.Stdout, "\tincrement major|minor|patch\t\tincrement a given version segment\n")
	fmt.Fprint(os.Stdout, "\tsatisfies <version constraint>\t\tverify a version satisfies a constraint\n")
	fmt.Fprint(os.Stdout, "\nusage:\n")
	fmt.Fprint(os.Stdout, "\tversion [command] [args] -f [filename]\n")
}

var filePath = flag.String("f", "", "the filename")

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
		os.Exit(1)
	}

	contents, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	command, args := flag.Args()[0], flag.Args()[1:]
	switch command {
	case "increment":
		increment(contents, args)
		os.Exit(0)
	case "satisfies":
		satisfies(contents, args)
		os.Exit(0)
	case "validate":
		validate(contents, args)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "invalid command: %s\n", command)
		os.Exit(1)
	}
}

func increment(contents []byte, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "must supply a semver segment\n")
		os.Exit(1)
	}

	version, err := semver.NewVersion(strings.TrimSpace(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	switch strings.ToLower(os.Args[1]) {
	case "patch":
		version.BumpPatch()
	case "minor":
		version.BumpMinor()
	case "major":
		version.BumpMajor()
	default:
		fmt.Fprintf(os.Stderr, "invalid segment: %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "should be one of: 'major', 'minor', and 'patch'\n")
		os.Exit(1)
	}

	fmt.Printf("%v", version)
}

func satisfies(contents []byte, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "must supply a version constraint\n")
		os.Exit(1)
	}

	version, err := semver.NewVersion(strings.TrimSpace(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	constraint, err := semver.NewConstraint(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if !constraint.Check(version) {
		fmt.Fprintf(os.Stderr, "`%s` does not satisfy `%s`\n", constraint.String(), version.String())
		os.Exit(1)
	}
}

func validate(contents []byte, args []string) {
	_, err := semver.NewVersion(strings.TrimSpace(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
