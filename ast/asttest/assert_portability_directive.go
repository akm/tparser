package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
)

// func AssertPortabilityDirectives(t *testing.T, expected, actual []ast.PortabilityDirective) {
// 	if !assert.Equal(t, len(expected), len(actual)) {
// 		return
// 	}
// 	for i, exp := range expected {
// 		act := actual[i]
// 		if !assert.Equal(t, exp, act) {
// 			AssertPortabilityDirective(t, exp, act)
// 		}
// 	}
// }

func AssertPortabilityDirective(t *testing.T, expected, actual *ast.PortabilityDirective) {
	// Do nothing
}
