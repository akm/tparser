package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func runType(t *testing.T, name string, text []rune, expected ast.Type, funcs ...func() interface{}) {
	t.Run(name, func(t *testing.T) {
		args := make([]interface{}, len(funcs))
		for i, f := range funcs {
			args[i] = f()
		}
		parser := NewTestParser(&text, args...)
		parser.NextToken()
		res, err := parser.ParseType()
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expected, res)
		}
	})
}

func runTypeDecl(t *testing.T, name string, text []rune, expected *ast.TypeDecl, funcs ...func() interface{}) {
	t.Run(name, func(t *testing.T) {
		args := make([]interface{}, len(funcs))
		for i, f := range funcs {
			args[i] = f()
		}
		parser := NewTestParser(&text, args...)
		parser.NextToken()
		res, err := parser.ParseTypeDecl()
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			if !assert.Equal(t, expected, res) {
				asttest.AssertTypeDecl(t, expected, res)
			}
		}
	})
}
