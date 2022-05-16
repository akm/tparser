package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertBlock(t *testing.T, expected, actual *ast.Block) {
	if !assert.Equal(t, expected.DeclSections, actual.DeclSections) {
		AssertDeclSections(t, expected.DeclSections, actual.DeclSections)
	}
	if !assert.Equal(t, expected.ExportsStmts1, actual.ExportsStmts1) {
		AssertExportsStmts(t, expected.ExportsStmts1, actual.ExportsStmts1)
	}
	if !assert.Equal(t, expected.Body, actual.Body) {
		AssertBlockBody(t, expected.Body, actual.Body)
	}
	if !assert.Equal(t, expected.ExportsStmts2, actual.ExportsStmts2) {
		AssertExportsStmts(t, expected.ExportsStmts2, actual.ExportsStmts2)
	}
}

func AssertExportsStmts(t *testing.T, expected, actual ast.ExportsStmts) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertExportsStmt(t, exp, act)
		}
	}
}

func AssertBlockBody(t *testing.T, expected, actual ast.BlockBody) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.CompoundStmt:
		AssertCompoundStmt(t, exp, actual.(*ast.CompoundStmt))
	case *ast.AssemblerStatement:
		AssertAssemblerStatement(t, exp, actual.(*ast.AssemblerStatement))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertExportsStmt(t *testing.T, expected, actual *ast.ExportsStmt) {
	if !assert.Equal(t, expected.ExportsItems, actual.ExportsItems) {
		AssertExportsItems(t, expected.ExportsItems, actual.ExportsItems)
	}
}

func AssertExportsItems(t *testing.T, expected, actual []*ast.ExportsItem) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertExportsItem(t, exp, act)
		}
	}
}

func AssertExportsItem(t *testing.T, expected, actual *ast.ExportsItem) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.Name, actual.Name) {
		AssertConstExpr(t, expected.Name, actual.Name)
	}
	if !assert.Equal(t, expected.Index, actual.Index) {
		AssertConstExpr(t, expected.Index, actual.Index)
	}
}

func AssertDeclSections(t *testing.T, expected, actual ast.DeclSections) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertDeclSection(t, exp, act)
		}
	}
}

func AssertDeclSection(t *testing.T, expected, actual ast.DeclSection) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.LabelDeclSection:
		AssertLabelDeclSection(t, exp, actual.(*ast.LabelDeclSection))
	case ast.ConstSection:
		AssertConstSection(t, exp, actual.(ast.ConstSection))
	case ast.TypeSection:
		AssertTypeSection(t, exp, actual.(ast.TypeSection))
	case ast.VarSection:
		AssertVarSection(t, exp, actual.(ast.VarSection))
	case *ast.FunctionDecl:
		AssertFunctionDecl(t, exp, actual.(*ast.FunctionDecl))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertLabelDeclSection(t *testing.T, expected, actual *ast.LabelDeclSection) {
	if !assert.Equal(t, expected.LabelId, actual.LabelId) {
		AssertIdent(t, expected.LabelId, actual.LabelId)
	}
}

func AssertLbelId(t *testing.T, expected, actual *ast.LabelId) {
	AssertIdent(t, expected, actual)
}
