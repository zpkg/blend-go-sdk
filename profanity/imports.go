package profanity

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/blend/go-sdk/exception"
)

// ImportsContainAny returns a profanity error if a given file contains any of a list of imports.
func ImportsContainAny(imports ...string) RuleFunc {
	return func(contents []byte) error {
		fset := token.NewFileSet()

		ast, err := parser.ParseFile(fset, "", contents, parser.ImportsOnly)
		if err != nil {
			return exception.New(err)
		}
		for _, fileImport := range ast.Imports {
			for _, i := range imports {
				if Glob(i, fileImport.Path.Value) {
					return fmt.Errorf("go import match: \"%s\"", i)
				}
			}
		}

		return nil
	}
}
