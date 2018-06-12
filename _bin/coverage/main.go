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
)

const (
	star = "*"
)

var reportOutputPath = flag.String("output", "coverage.html", "the path to write the full html coverage report")
var temporaryOutputPath = flag.String("tmp", "coverage.cov", "the path to write the intermediate results")
var update = flag.Bool("update", false, "if we should write the current coverage to `COVERAGE` files")
var enforce = flag.Bool("enforce", false, "if we should enforce coverage minimums defined in `COVERAGE` files")
var include = flag.String("include", "", "the include file filter in glob form, can be a csv.")
var exclude = flag.String("exclude", "", "the exclude file filter in glob form, can be a csv.")

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

		if info.Name() == "vendor" {
			return filepath.SkipDir
		}

		if !dirHasGlob(currentPath, "*.go") {
			return nil
		}

		if len(*include) > 0 {
			if matches := globAnyMatch(*include, currentPath); !matches {
				return nil
			}
		}

		if len(*exclude) > 0 {
			if matches := globAnyMatch(*exclude, currentPath); matches {
				return nil
			}
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
			return err
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

// globIncludeMatch tests if a file matches a (potentially) csv of glob filters.
func globAnyMatch(filter, file string) bool {
	parts := strings.Split(filter, ",")
	for _, part := range parts {
		if matches := glob(strings.TrimSpace(part), file); matches {
			return true
		}
	}
	return false
}

func glob(pattern, subj string) bool {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == star {
		return true
	}

	parts := strings.Split(pattern, star)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, star)
	trailingGlob := strings.HasSuffix(pattern, star)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
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
	return ioutil.WriteFile(filepath.Join(path, "COVERAGE"), []byte(strings.TrimSpace(coverage)), 0755)
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
