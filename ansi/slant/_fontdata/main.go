package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	contents, err := ioutil.ReadFile("slant_letters.txt")
	if err != nil {
		log.Fatal(err)
	}

	letterHeight := 6

	scanner := bufio.NewScanner(bytes.NewReader(contents))

	var lineText string
	var line int

	fmt.Println("[][]string{")

	for scanner.Scan() {
		lineText = scanner.Text()

		if line == 0 {
			fmt.Println("\t{")
		}
		if line < letterHeight {
			lineText = strings.TrimSuffix(strings.TrimSuffix(lineText, "@"), "@")
			lineText = strings.ReplaceAll(lineText, "\"", "\\\"") // escape quotes
			lineText = strings.ReplaceAll(lineText, "\\", "\\\\") // escape slashes
			fmt.Printf("\t\t\"%s\",\n", lineText)
			line++
		}

		if line == letterHeight || strings.HasSuffix(lineText, "@@") {
			fmt.Println("\t},")
			line = 0
		}
	}

	fmt.Println("}")
}
