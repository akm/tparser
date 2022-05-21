package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser/parsertest"
	"github.com/stretchr/testify/assert"
)

func runProgram(t *testing.T, name string, clearLocations bool, text []rune, expected *ast.Program, callbacks ...func(expected, actual *ast.CompoundStmt)) {
	t.Run(name, func(t *testing.T) {
		parser := parsertest.NewTestParser(&text, parsertest.NewTestProgramContext())
		parser.NextToken()
		res, err := parser.ParseProgram()
		if assert.NoError(t, err) {
			if clearLocations {
				asttest.ClearLocations(t, res)
			}
			if !assert.Equal(t, expected, res) {
				asttest.AssertProgram(t, expected, res)
			}
		}
	})
}

func runSatement(t *testing.T, name string, clearLocations bool, text []rune, expected *ast.Statement, callbacks ...func(expected, actual *ast.Statement)) {
	t.Run(name, func(t *testing.T) {
		parser := parsertest.NewTestParser(&text)
		parser.NextToken()
		res, err := parser.ParseStatement()
		if assert.NoError(t, err) {
			if clearLocations {
				asttest.ClearLocations(t, res)
			}
			if !assert.Equal(t, expected, res) {
				for _, callback := range callbacks {
					callback(expected, res)
				}
			}
		}
	})
}
