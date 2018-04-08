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
	"regexp"
	"strconv"
	"strings"

	"github.com/blend/go-sdk/exception"
)

var reportOutputPath = flag.String("output", "coverage.html", "the path to write the full html coverage report")
var temporaryOutputPath = flag.String("tmp", "coverage.cov", "the path to write the intermediate results")
var update = flag.Bool("update", false, "if we should write the current coverage to `COVERAGE` files")
var enforce = flag.Bool("enforce", false, "if we should enforce coverage minimums defined in `COVERAGE` files")

func main() {
	flag.Parse()

	fmt.Fprintln(os.Stdout, "coverage starting")
	tempOutput, err := removeAndOpen(*temporaryOutputPath)
	if err != nil {
		maybeFatal(err)
	}
	fmt.Fprintln(tempOutput, "mode: set")

	maybeFatal(filepath.Walk("./", func(currentPath string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
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

		intermediateFile := filepath.Join(currentPath, "profile.cov")
		err = removeIfExists(intermediateFile)
		if err != nil {
			return err
		}

		var output []byte
		output, err = execCoverage(currentPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, string(output))
			return exception.Wrap(err)
		}
		coverage := extractCoverage(string(output))
		fmt.Fprintf(os.Stdout, "%s: %v%%\n", currentPath, coverage)

		if enforce != nil && *enforce {
			err = enforceCoverage(currentPath, coverage)
			if err != nil {
				return err
			}
		}

		if update != nil && *update {
			fmt.Fprintf(os.Stdout, "%s updating coverage\n", currentPath)
			err = writeCoverage(currentPath, coverage)
			if err != nil {
				return err
			}
		}

		err = mergeCoverageOutput(intermediateFile, tempOutput)
		if err != nil {
			return err
		}

		err = removeIfExists(intermediateFile)
		if err != nil {
			return err
		}

		return nil
	}))

	maybeFatal(tempOutput.Close())

	fmt.Fprintf(os.Stdout, "merging coverage output: %s\n", *reportOutputPath)
	maybeFatal(execCoverageCompile())
	maybeFatal(removeIfExists(*temporaryOutputPath))
	fmt.Fprintln(os.Stdout, "coverage complete")
}

func enforceCoverage(path, actualCoverage string) error {
	actual, err := strconv.ParseFloat(actualCoverage, 64)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(filepath.Join(path, "COVERAGE"))
	if err != nil {
		return err
	}
	expected, err := strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		return err
	}

	if expected == 0 {
		return nil
	}

	if actual < expected {
		return fmt.Errorf("%s fails coverage: %0.2f%% vs. %0.2f%%", path, expected, actual)
	}
	return nil
}

func extractCoverage(corpus string) string {
	regex := `coverage: ([0-9,.]+)% of statements`
	expr := regexp.MustCompile(regex)

	results := expr.FindStringSubmatch(corpus)
	if len(results) > 1 {
		return results[1]
	}
	return "0"
}

func writeCoverage(path, coverage string) error {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(coverage), 64)
	if err != nil {
		return err
	}

	expected := strconv.Itoa(int(parsed))
	return ioutil.WriteFile(filepath.Join(path, "COVERAGE"), []byte(expected), 0755)
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
	return cmd.Run()
}

func mergeCoverageOutput(temp string, outFile *os.File) error {
	contents, err := ioutil.ReadFile(temp)
	if err != nil {
		return err
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
			return err
		}
	}
	return nil
}

func removeIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
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
			return nil, err
		}
	}
	return os.Create(path)
}
