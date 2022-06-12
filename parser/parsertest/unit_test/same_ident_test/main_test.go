package sameident_test

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestSameIdent(t *testing.T) {
	actualProject1, err := parser.ParseProgram("Project1.dpr")
	if !assert.NoError(t, err) {
		return
	}
	asttest.ClearLocations(t, actualProject1.Program)
	asttest.ClearUnitDeclMaps(t, actualProject1.Program)
	actualUnit1 := actualProject1.Units.ByName("Unit1")
	actualUnit2 := actualProject1.Units.ByName("Unit2")
	actualUnit3 := actualProject1.Units.ByName("Unit3")
	actualUnit4 := actualProject1.Units.ByName("Unit4")

	actualUnits := ast.Units{actualUnit1, actualUnit2, actualUnit3, actualUnit4}
	for _, unit := range actualUnits {
		if !assert.NotNil(t, unit) {
			return
		}
	}

	for _, unit := range actualUnits {
		asttest.ClearLocations(t, unit)
		unit.DeclMap = nil
	}

	newStringVarDecl := func(name, value string) *ast.VarDecl {
		return &ast.VarDecl{
			IdentList: asttest.NewIdentList(asttest.NewIdent(name)),
			Type:      asttest.NewStringType("string"),
			ConstExpr: asttest.NewConstExpr(asttest.NewString("'" + value + "'")),
		}
	}

	writeln := func(arg interface{}) *ast.Statement {
		var writelnArg *ast.Expression
		switch v := arg.(type) {
		case *ast.Expression:
			writelnArg = v
		case *ast.IdentRef:
			writelnArg = asttest.NewExpression(v)
		case *ast.VarDecl:
			identRef := asttest.NewIdentRef(v.IdentList[0].Name, v.ToDeclarations()[0])
			writelnArg = asttest.NewExpression(identRef)
		case *ast.QualId:
			writelnArg = asttest.NewExpression(v)
		case string:
			writelnArg = asttest.NewExpression(asttest.NewString("'" + v + "'"))
		default:
			t.Fatalf("unexpected type: %T", arg)
		}

		return &ast.Statement{
			Body: &ast.CallStatement{
				Designator: asttest.NewDesignator(asttest.NewIdentRef("Writeln")),
				ExprList:   ast.ExprList{writelnArg},
			},
		}
	}
	readln := func() *ast.Statement {
		return &ast.Statement{
			Body: &ast.CallStatement{
				Designator: asttest.NewDesignator(asttest.NewIdentRef("Readln")),
			},
		}
	}

	expectUnit1Foo := newStringVarDecl("Foo", "Foo@Unit1")
	expectUnit1Bar := newStringVarDecl("Bar", "Bar@Unit1")
	expectUnit1Unit2 := newStringVarDecl("Unit2", "Unit2@Unit1")
	expectUnit1Project1 := newStringVarDecl("Project1", "Project1@Unit1")

	expectUnit2Foo := newStringVarDecl("Foo", "Foo@Unit2")
	expectUnit2Bar := newStringVarDecl("Bar", "Bar@Unit2")
	expectUnit2Unit1 := newStringVarDecl("Unit1", "Unit1@Unit2")
	expectUnit2Project1 := newStringVarDecl("Project1", "Project1@Unit2")

	expectUnit3Baz := newStringVarDecl("Baz", "Baz@Unit3") // in implementation section

	expectUnit4Foo := newStringVarDecl("Foo", "Foo@Unit4")
	expectUnit4Bar := newStringVarDecl("Bar", "Bar@Unit4")
	expectUnit4Project1 := newStringVarDecl("Project1", "Project1@Unit4")

	expectProject1Bar := newStringVarDecl("Bar", "Bar@Project1")
	expectProject1Baz := newStringVarDecl("Baz", "Baz@Project1")

	expectUnit1Unit4 := asttest.NewUnitRef("Unit4")

	expectUnit1 := &ast.Unit{
		Path:  "Unit1.pas",
		Ident: asttest.NewIdent("Unit1"),
		InterfaceSection: &ast.InterfaceSection{
			UsesClause: ast.UsesClause{
				asttest.NewUnitRef("SysUtils"),
				expectUnit1Unit4,
			},
			InterfaceDecls: ast.InterfaceDecls{
				ast.VarSection{
					expectUnit1Foo, expectUnit1Bar, expectUnit1Unit2, expectUnit1Project1,
				},
			},
		},
		ImplementationSection: &ast.ImplementationSection{},
	}

	expectUnit1.InitSection = &ast.InitSection{
		InitializationStmts: ast.StmtList{
			writeln("-- Unit1 Initialization----"),
			writeln(expectUnit1Foo),
			writeln(asttest.NewQualId(
				asttest.NewIdentRef("Unit1", expectUnit1.ToDeclarations()[0]),
				asttest.NewIdentRef("Foo", expectUnit1Foo.ToDeclarations()[0]),
			)),
			writeln(expectUnit1Bar),
			writeln(asttest.NewQualId(
				asttest.NewIdentRef("Unit4", expectUnit1Unit4.ToDeclarations()[0]),
				asttest.NewIdentRef("Foo", expectUnit4Foo.ToDeclarations()[0]),
			)),
			writeln(asttest.NewQualId(
				asttest.NewIdentRef("Unit4", expectUnit1Unit4.ToDeclarations()[0]),
				asttest.NewIdentRef("Bar", expectUnit4Bar.ToDeclarations()[0]),
			)),
			writeln(expectUnit1Unit2),
		},
		FinalizationStmts: ast.StmtList{
			writeln("-- Unit1 Finalization----"),
			writeln(expectUnit1Foo),
			writeln(expectUnit1Bar),
			writeln(expectUnit1Unit2),
			readln(),
		},
	}

	expectUnit2Unit1InImpl := asttest.NewUnitRef("Unit1")
	expectUnit2Unit1InImpl.Unit = expectUnit1

	expectUnit2 := &ast.Unit{
		Path:  "Unit2.pas",
		Ident: asttest.NewIdent("Unit2"),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{
				ast.VarSection{
					expectUnit2Foo, expectUnit2Bar, expectUnit2Unit1, expectUnit2Project1,
				},
			},
		},
		ImplementationSection: &ast.ImplementationSection{
			UsesClause: ast.UsesClause{expectUnit2Unit1InImpl},
		},
		InitSection: &ast.InitSection{
			InitializationStmts: ast.StmtList{
				writeln("-- Unit2 Initialization----"),
				writeln(expectUnit2Foo),
				writeln(expectUnit2Bar),
				// writeln(expectUnit2Unit1),
				writeln(asttest.NewQualId(
					asttest.NewIdentRef("Unit1", expectUnit2Unit1InImpl.ToDeclarations()[0]),
					asttest.NewIdentRef("Foo", expectUnit1Foo.ToDeclarations()[0]),
				)),
				writeln(asttest.NewQualId(
					asttest.NewIdentRef("Unit1", expectUnit2Unit1InImpl.ToDeclarations()[0]),
					asttest.NewIdentRef("Bar", expectUnit1Bar.ToDeclarations()[0]),
				)),
			},
			FinalizationStmts: ast.StmtList{
				writeln("-- Unit2 Finalization----"),
				writeln(expectUnit2Foo),
				writeln(expectUnit2Bar),
				// writeln(expectUnit2Unit1),
				writeln(asttest.NewQualId(
					asttest.NewIdentRef("Unit1", expectUnit2Unit1InImpl.ToDeclarations()[0]),
					asttest.NewIdentRef("Foo", expectUnit1Foo.ToDeclarations()[0]),
				)),
				writeln(asttest.NewQualId(
					asttest.NewIdentRef("Unit1", expectUnit2Unit1InImpl.ToDeclarations()[0]),
					asttest.NewIdentRef("Bar", expectUnit1Bar.ToDeclarations()[0]),
				)),
				readln(),
			},
		},
	}

	expectUnit3 := &ast.Unit{
		Path:  "Unit3.pas",
		Ident: asttest.NewIdent("Unit3"),
		InterfaceSection: &ast.InterfaceSection{
			UsesClause: ast.UsesClause{
				asttest.NewUnitRef("Unit1"),
				asttest.NewUnitRef("Unit2"),
			},
		},
		ImplementationSection: &ast.ImplementationSection{
			DeclSections: ast.DeclSections{
				ast.VarSection{expectUnit3Baz},
			},
		},
	}
	expectUnit3.InitSection = &ast.InitSection{
		InitializationStmts: ast.StmtList{
			writeln("-- Unit3 Initialization----"),
			writeln(expectUnit2Foo),
			writeln(expectUnit2Bar),
			writeln(expectUnit3Baz),
			writeln(asttest.NewQualId(
				asttest.NewIdentRef("Unit3", expectUnit3.ToDeclarations()[0]),
				asttest.NewIdentRef("Baz", expectUnit3Baz.ToDeclarations()[0]),
			)),
			writeln(expectUnit2Project1),
		},
		FinalizationStmts: ast.StmtList{
			writeln("-- Unit3 Finalization----"),
			writeln(expectUnit2Foo),
			writeln(expectUnit2Bar),
			writeln(expectUnit2Project1),
			readln(),
		},
	}

	expectUnit4 := &ast.Unit{
		Path:  "Unit4.pas",
		Ident: asttest.NewIdent("Unit4"),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{
				ast.VarSection{
					expectUnit4Foo, expectUnit4Bar, expectUnit4Project1,
				},
			},
		},
		ImplementationSection: &ast.ImplementationSection{},
		InitSection: &ast.InitSection{
			InitializationStmts: ast.StmtList{
				writeln("-- Unit4 Initialization----"),
				writeln(expectUnit4Foo),
				writeln(expectUnit4Bar),
				writeln(expectUnit4Project1),
			},
			FinalizationStmts: ast.StmtList{
				writeln("-- Unit4 Finalization----"),
				writeln(expectUnit4Foo),
				writeln(expectUnit4Bar),
				writeln(expectUnit4Project1),
				readln(),
			},
		},
	}

	expectedProject1 := &parser.Program{
		Program: &ast.Program{
			Path:  "Project1.dpr",
			Ident: asttest.NewIdent("Project1"),
		},
		Units: ast.Units{expectUnit1, expectUnit2, expectUnit3, expectUnit4},
	}
	project1Unit1 := &ast.UsesClauseItem{Ident: asttest.NewIdent("Unit1"), Path: ext.StringPtr("'Unit1.pas'"), Unit: expectUnit1}
	project1Unit2 := &ast.UsesClauseItem{Ident: asttest.NewIdent("Unit2"), Path: ext.StringPtr("'Unit2.pas'"), Unit: expectUnit2}
	project1Unit3 := &ast.UsesClauseItem{Ident: asttest.NewIdent("Unit3"), Path: ext.StringPtr("'Unit3.pas'"), Unit: expectUnit3}
	project1Unit4 := &ast.UsesClauseItem{Ident: asttest.NewIdent("Unit4"), Path: ext.StringPtr("'Unit4.pas'"), Unit: expectUnit4}
	expectedProject1.Program.ProgramBlock = &ast.ProgramBlock{
		UsesClause: ast.UsesClause{
			{Ident: asttest.NewIdent("SysUtils")},
			project1Unit1,
			project1Unit2,
			project1Unit3,
			project1Unit4,
		},
		Block: &ast.Block{
			DeclSections: ast.DeclSections{
				ast.VarSection{expectProject1Bar, expectProject1Baz},
			},
			Body: &ast.CompoundStmt{
				StmtList: ast.StmtList{
					writeln("-- Project1 ----"),
					writeln(expectUnit4Foo),
					writeln(expectProject1Bar),
					writeln(asttest.NewQualId(
						asttest.NewIdentRef("Unit1", project1Unit1.ToDeclarations()[0]),
						asttest.NewIdentRef("Unit2", expectUnit1Unit2.ToDeclarations()[0]),
					)),
					writeln(asttest.NewQualId(
						asttest.NewIdentRef("Unit2", project1Unit2.ToDeclarations()[0]),
						asttest.NewIdentRef("Unit1", expectUnit2Unit1.ToDeclarations()[0]),
					)),
					writeln(asttest.NewQualId(
						asttest.NewIdentRef("Project1", expectedProject1.Program.ToDeclarations()[0]),
						asttest.NewIdentRef("Bar", expectProject1Bar.ToDeclarations()[0]),
					)),
					writeln(asttest.NewQualId(
						asttest.NewIdentRef("Project1", expectedProject1.Program.ToDeclarations()[0]),
						asttest.NewIdentRef("Baz", expectProject1Baz.ToDeclarations()[0]),
					)),
					readln(),
				},
			},
		},
	}

	func() {
		for _, item := range expectUnit1.InterfaceSection.UsesClause {
			if item.Ident.Name == expectUnit4.Ident.Name {
				item.Unit = expectUnit4
			}
		}
		for _, item := range expectUnit3.InterfaceSection.UsesClause {
			if item.Ident.Name == expectUnit1.Ident.Name {
				item.Unit = expectUnit1
			} else if item.Ident.Name == expectUnit2.Ident.Name {
				item.Unit = expectUnit2
			}
		}
	}()

	t.Run("Unit1", func(t *testing.T) {
		if !assert.Equal(t, expectUnit1, actualUnit1) {
			asttest.AssertUnit(t, expectUnit1, actualUnit1)
		}
	})
	t.Run("Unit2", func(t *testing.T) {
		if !assert.Equal(t, expectUnit2, actualUnit2) {
			asttest.AssertUnit(t, expectUnit2, actualUnit2)
		}
	})
	t.Run("Unit3", func(t *testing.T) {
		if !assert.Equal(t, expectUnit3, actualUnit3) {
			asttest.AssertUnit(t, expectUnit3, actualUnit3)
		}
	})
	t.Run("Unit4", func(t *testing.T) {
		if !assert.Equal(t, expectUnit4, actualUnit4) {
			asttest.AssertUnit(t, expectUnit4, actualUnit4)
		}
	})
	t.Run("Project1", func(t *testing.T) {
		actualProject1.DeclMap = nil
		if !assert.Equal(t, expectedProject1, actualProject1) {
			asttest.AssertProgram(t, expectedProject1.Program, actualProject1.Program)
		}
	})
}
