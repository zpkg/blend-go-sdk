package profanity

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

// ImportsContainAny returns a profanity error if a given file contains any of a list of imports.
func ImportsContainAny(imports ...string) RuleFunc {
	return func(filename string, contents []byte) RuleResult {
		fset := token.NewFileSet()

		ast, err := parser.ParseFile(fset, filename, contents, parser.ImportsOnly)
		if err != nil {
			return RuleResult{Err: err}
		}
		for _, fileImport := range ast.Imports {
			for _, i := range imports {
				if Glob(i, strings.Trim(fileImport.Path.Value, "\"")) {
					return RuleResult{
						File:    filename,
						Line:    fset.Position(fileImport.Pos()).Line,
						Message: fmt.Sprintf("go imports include: \"%s\"", i),
					}
				}
			}
		}
		return RuleResult{OK: true}
	}
}
