package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Rules is the full rule suite that runs for every file.
	Rules = []Rule{
		Contains("github.com/blendlabs/"),
		Contains("gopkg.in/"),
	}
)

func main() {
	// walk the filesystem
	// for each file named by the gob filter
	// run the rules on it

	walkErr := filepath.Walk("./", func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasSuffix(file, ".git") {
			return filepath.SkipDir
		}
		if info.IsDir() && strings.HasSuffix(file, "_bin") {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(file, ".go") {
			return nil
		}

		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		for _, rule := range Rules {
			if err := rule(contents); err != nil {
				return fmt.Errorf("%s failed: %+v", file, err)
			}
		}

		return nil
	})

	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", walkErr)
		os.Exit(1)
		return
	}
	fmt.Fprintf(os.Stdout, "profanity ok!\n")
	os.Exit(0)
}

// Contains creates a simple contains rule.
func Contains(value string) Rule {
	return func(contents []byte) error {
		if strings.Contains(string(contents), value) {
			return fmt.Errorf("contains: \"%s\"", value)
		}
		return nil
	}
}

// Regex creates a new regex filter rule.
func Regex(expr string) Rule {
	regex := regexp.MustCompile(expr)
	return func(contents []byte) error {
		if regex.Match(contents) {
			return fmt.Errorf("regexp match: \"%s\"", expr)
		}
		return nil
	}
}

// Rule evaluates contents.
type Rule func([]byte) error
