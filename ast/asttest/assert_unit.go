package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertUnits(t *testing.T, expected, actual ast.Units) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertUnit(t, exp, act)
		}
	}
}

func AssertUnit(t *testing.T, expected, actual *ast.Unit) {
	assert.Equal(t, expected.Path, actual.Path)
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.PortabilityDirective, actual.PortabilityDirective) {
		AssertPortabilityDirective(t, expected.PortabilityDirective, actual.PortabilityDirective)
	}
	if !assert.Equal(t, expected.InterfaceSection, actual.InterfaceSection) {
		AssertInterfaceSection(t, expected.InterfaceSection, actual.InterfaceSection)
	}
	if !assert.Equal(t, expected.ImplementationSection, actual.ImplementationSection) {
		AssertImplementationSection(t, expected.ImplementationSection, actual.ImplementationSection)
	}
	if !assert.Equal(t, expected.InitSection, actual.InitSection) {
		AssertInitSection(t, expected.InitSection, actual.InitSection)
	}
}

func AssertInterfaceSection(t *testing.T, expected, actual *ast.InterfaceSection) {
	if !assert.Equal(t, expected.UsesClause, actual.UsesClause) {
		AssertUsesClause(t, expected.UsesClause, actual.UsesClause)
	}
	if !assert.Equal(t, expected.InterfaceDecls, actual.InterfaceDecls) {
		AssertInterfaceDecls(t, expected.InterfaceDecls, actual.InterfaceDecls)
	}
}

func AssertInterfaceDecls(t *testing.T, expected, actual ast.InterfaceDecls) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertInterfaceDecl(t, exp, act)
		}
	}
}

func AssertInterfaceDecl(t *testing.T, expected, actual ast.InterfaceDecl) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case ast.ConstSection:
		AssertConstSection(t, exp, actual.(ast.ConstSection))
	case ast.TypeSection:
		AssertTypeSection(t, exp, actual.(ast.TypeSection))
	case ast.VarSection:
		AssertVarSection(t, exp, actual.(ast.VarSection))
	case ast.ThreadVarSection:
		AssertThreadVarSection(t, exp, actual.(ast.ThreadVarSection))
	case *ast.ExportedHeading:
		AssertExportedHeading(t, exp, actual.(*ast.ExportedHeading))
	default:
		assert.Fail(t, "unexpected interface decl type: %T", exp)
	}
}

func AssertImplementationSection(t *testing.T, expected, actual *ast.ImplementationSection) {
	if !assert.Equal(t, expected.UsesClause, actual.UsesClause) {
		AssertUsesClause(t, expected.UsesClause, actual.UsesClause)
	}
	if !assert.Equal(t, expected.DeclSections, actual.DeclSections) {
		AssertDeclSections(t, expected.DeclSections, actual.DeclSections)
	}
	if !assert.Equal(t, expected.ExportsStmts, actual.ExportsStmts) {
		AssertExportsStmts(t, expected.ExportsStmts, actual.ExportsStmts)
	}
}

func AssertInitSection(t *testing.T, expected, actual *ast.InitSection) {
	if !assert.Equal(t, expected.InitializationStmts, actual.InitializationStmts) {
		AssertStmtList(t, expected.InitializationStmts, actual.InitializationStmts)
	}
	if !assert.Equal(t, expected.FinalizationStmts, actual.FinalizationStmts) {
		AssertStmtList(t, expected.FinalizationStmts, actual.FinalizationStmts)
	}
}

func AssertUnitId(t *testing.T, expected, actual *ast.UnitId) {
	AssertIdent(t, (*ast.Ident)(expected), (*ast.Ident)(actual))
}

func AssertQualId(t *testing.T, expected, actual *ast.QualId) {
	if !assert.Equal(t, expected.UnitId, actual.UnitId) {
		AssertIdentRef(t, expected.UnitId, actual.UnitId)
	}
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdentRef(t, expected.Ident, actual.Ident)
	}
}

func AssertQualIds(t *testing.T, expected, actual ast.QualIds) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertQualId(t, exp, act)
		}
	}
}
