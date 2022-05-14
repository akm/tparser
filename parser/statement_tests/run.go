package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func runProgram(t *testing.T, name string, clearLocations bool, text []rune, expected *ast.Program, callbacks ...func(expected, actual *ast.CompoundStmt)) {
	t.Run(name, func(t *testing.T) {
		parser := parser.NewParser(&text)
		parser.NextToken()
		res, err := parser.ParseProgram()
		if assert.NoError(t, err) {
			if clearLocations {
				asttest.ClearLocations(t, res)
			}
			if !assert.Equal(t, expected, res) {
				assert.Equal(t, expected.ProgramBlock.UsesClause, res.ProgramBlock.UsesClause)
				if !assert.Equal(t, expected.ProgramBlock.Block, res.ProgramBlock.Block) {
					assert.Equal(t, expected.ProgramBlock.Block.DeclSections, res.ProgramBlock.Block.DeclSections)
					assert.Equal(t, expected.ProgramBlock.Block.ExportsStmts1, res.ProgramBlock.Block.ExportsStmts1)
					if !assert.Equal(t, expected.ProgramBlock.Block.CompoundStmt, res.ProgramBlock.Block.CompoundStmt) {
					}
					assert.Equal(t, expected.ProgramBlock.Block.ExportsStmts2, res.ProgramBlock.Block.ExportsStmts2)
				}
			}
		}
	})
}

func runSatement(t *testing.T, name string, clearLocations bool, text []rune, expected *ast.Statement, callbacks ...func(expected, actual *ast.Statement)) {
	t.Run(name, func(t *testing.T) {
		parser := parser.NewParser(&text)
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
