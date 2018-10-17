package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blend/go-sdk/yaml"
)

type configVar struct {
	Field   string
	Default string
	Value   string
}

func (cv *configVar) Set(value string) error {
	cv.Value = value
	return nil
}

type configVarSet []configVar

func (cvs configVarSet) String() string {
	return "Fields to prompt for in the form fieldname=default"
}

func (cvs *configVarSet) Set(flagValue string) error {
	parts := strings.SplitN(flagValue, "=", 2)
	if len(parts) > 1 {
		*cvs = append(*cvs, configVar{
			Field:   parts[0],
			Default: parts[1],
		})
	} else if len(parts) > 0 {
		*cvs = append(*cvs, configVar{
			Field: parts[0],
		})
	} else {
		return fmt.Errorf("invalid config var")
	}
	return nil
}

func (cvs configVarSet) MarshalYAML() (interface{}, error) {
	output := make(map[string]interface{})
	for _, cv := range cvs {
		if len(cv.Value) > 0 {
			output[cv.Field] = cv.Value
		} else {
			output[cv.Field] = cv.Default
		}
	}
	return output, nil
}

func main() {
	fields := configVarSet{}
	output := flag.String("o", "vars.yml", "The output file")
	flag.Var(&fields, "field", "Fields to prompt for, can be multiple")
	flag.Parse()

	if len(fields) == 0 {
		flag.Usage()
	}

	for _, field := range fields {
		scanln(&field)
	}

	buffer := new(bytes.Buffer)
	if err := yaml.NewEncoder(buffer).Encode(fields); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(*output, buffer.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func scanln(cv *configVar) {
	if len(cv.Default) > 0 {
		fmt.Printf("%s[%s]: ", cv.Field, cv.Default)
	} else {
		fmt.Printf("%s: ", cv.Field)
	}
	fmt.Scanln(&cv.Value)
}
