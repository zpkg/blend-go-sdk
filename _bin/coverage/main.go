package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/exception"
)

var reportOutputPath = flag.String("output", "coverage.html", "the path to write the full html coverage report")
var temporaryOutputPath = flag.String("tmp", "coverage.cov", "the path to write the intermediate results")
var enforce = flag.Bool("enforce", false, "if we should enforce coverage minimums defined in `COVERAGE` files")

func main() {
	tempOutput, err := removeAndOpen(*temporaryOutputPath)
	if err != nil {
		maybeFatal(err)
	}
	fmt.Fprintln(tempOutput, "mode: set")

	maybeFatal(filepath.Walk("./", func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if info.Name() == ".git" {
			return filepath.SkipDir
		}

		if strings.HasPrefix(info.Name(), "_") {
			return filepath.SkipDir
		}

		if !dirHasGlob(currentPath, "*.go") {
			return nil
		}

		fmt.Fprintf(os.Stdout, "running coverage for: %s\n", currentPath)

		intermediateFile := filepath.Join(currentPath, "profile.cov")

		err = removeIfExists(intermediateFile)
		if err != nil {
			return err
		}

		var output []byte
		output, err = execCoverage(currentPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, string(output))
			return exception.Wrap(err)
		}

		err = mergeCoverageOutput(intermediateFile, tempOutput)
		if err != nil {
			return err
		}

		err = removeIfExists(intermediateFile)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "running coverage for: %s complete\n", currentPath)
		return nil
	}))

	maybeFatal(tempOutput.Close())
	maybeFatal(execCoverageCompile())
}

func dirHasGlob(path, glob string) bool {
	files, _ := filepath.Glob(filepath.Join(path, glob))
	return len(files) > 0
}

func gobin() string {
	gobin, err := exec.LookPath("go")
	maybeFatal(err)
	return gobin
}

func execCoverage(path string) ([]byte, error) {
	cmd := exec.Command(gobin(), "test", "-short", "-covermode=set", "-coverprofile=profile.cov")
	cmd.Dir = path
	return cmd.CombinedOutput()
}

func execCoverageCompile() error {
	cmd := exec.Command(gobin(), "tool", "cover", fmt.Sprintf("-html=%s", *temporaryOutputPath), fmt.Sprintf("-o=%s", *reportOutputPath))
	return exception.Wrap(cmd.Run())
}

func mergeCoverageOutput(temp string, outFile *os.File) error {
	contents, err := ioutil.ReadFile(temp)
	if err != nil {
		return exception.Wrap(err)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(contents))

	var skip int
	for scanner.Scan() {
		skip++
		if skip < 2 {
			continue
		}
		_, err = fmt.Fprintln(outFile, scanner.Text())
		if err != nil {
			return exception.Wrap(err)
		}
	}
	return nil
}

func removeIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return exception.Wrap(os.Remove(path))
	}
	return nil
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func removeAndOpen(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return nil, exception.Wrap(err)
		}
	}
	return os.Create(path)
}
