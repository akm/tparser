package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertProgram(t *testing.T, expected *ast.Program, actual *ast.Program) {
	assert.Equal(t, expected.Path, actual.Path)
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.ProgramBlock, actual.ProgramBlock) {
		AssertProgramBlock(t, expected.ProgramBlock, actual.ProgramBlock)
	}
}

func AssertProgramBlock(t *testing.T, expected *ast.ProgramBlock, actual *ast.ProgramBlock) {
	assert.Equal(t, expected.UsesClause, actual.UsesClause)
	if !assert.Equal(t, expected.Block, actual.Block) {
		AssertBlock(t, expected.Block, actual.Block)
	}
}
