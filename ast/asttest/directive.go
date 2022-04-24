package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewExternalOptions(libraryName string, args ...interface{}) *ast.ExternalOptions {
	return ast.NewExternalOptions(libraryName, args...)
}
