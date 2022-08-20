package methodcalling_test

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestMethodCalling(t *testing.T) {
	actualProject1, err := parser.ParseProgram("Project1.dpr")
	if !assert.NoError(t, err) {
		return
	}
	asttest.ClearLocations(t, actualProject1.Program)
	asttest.ClearUnitDeclMaps(t, actualProject1.Program)

	actualUnit1 := actualProject1.Units.ByName("Unit1")
	if !assert.NotNil(t, actualUnit1) {
		return
	}
	asttest.ClearLocations(t, actualUnit1)
	actualUnit1.DeclMap = nil

	declClass := func(name string, ancestorDecl *ast.TypeDecl, directives ast.ClassMethodDirectiveList) *ast.TypeDecl {
		var heritage ast.ClassHeritage
		if ancestorDecl != nil {
			heritage = ast.ClassHeritage{
				asttest.NewTypeId(ancestorDecl.Name, ancestorDecl.ToDeclarations()[0]),
			}
		}

		return &ast.TypeDecl{
			Ident: asttest.NewIdent(name),
			Type: &ast.CustomClassType{
				Heritage: heritage,
				Members: ast.ClassMemberSections{
					&ast.ClassMemberSection{
						Visibility: ast.CvPublic,
						ClassMethodList: ast.ClassMethodList{
							&ast.ClassMethod{
								Heading: &ast.FunctionHeading{
									Type:  ast.FtProcedure,
									Ident: asttest.NewIdent("Foo"),
								},
								Directives: directives,
							},
						},
					},
				},
			},
		}
	}

	classDeclT0 := declClass("T0", nil, ast.ClassMethodDirectiveList{ast.CmdVirtual, ast.CmdAbstract})
	classDeclT1 := declClass("T1", classDeclT0, ast.ClassMethodDirectiveList{ast.CmdOverride})
	classDeclT2 := declClass("T2", classDeclT1, ast.ClassMethodDirectiveList{ast.CmdOverride})
	classDeclT3 := declClass("T3", classDeclT1, ast.ClassMethodDirectiveList{ast.CmdReintroduce})
	classDeclT4 := declClass("T4", classDeclT1, ast.ClassMethodDirectiveList{})

	declFooMethod := func(classDecl *ast.TypeDecl) *ast.FunctionDecl {
		identRef := asttest.NewIdentRef(classDecl.Name, classDecl.ToDeclarations()[0])
		return &ast.FunctionDecl{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Class: identRef,
				Ident: asttest.NewIdent("Foo"),
			},
			Block: &ast.Block{
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						{
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator(
									asttest.NewQualId(
										asttest.NewIdent("Writeln"),
									),
								),
								ExprList: ast.ExprList{
									asttest.NewExpression(
										asttest.NewString(classDecl.Name + ".Foo"),
									),
								},
							},
						},
					},
				},
			},
		}
	}

	t1Foo := declFooMethod(classDeclT1)
	t2Foo := declFooMethod(classDeclT2)
	t3Foo := declFooMethod(classDeclT3)
	t4Foo := declFooMethod(classDeclT4)

	expectUnit1 := &ast.Unit{
		Path:  "Unit1.pas",
		Ident: asttest.NewIdent("Unit1"),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{
				ast.TypeSection{
					classDeclT0,
					classDeclT1,
					classDeclT2,
					classDeclT3,
					classDeclT4,
				},
			},
		},
		ImplementationSection: &ast.ImplementationSection{
			DeclSections: ast.DeclSections{
				t1Foo,
				t2Foo,
				t3Foo,
				t4Foo,
			},
		},
	}

	varDeclT0s := &ast.VarDecl{
		IdentList: asttest.NewIdentList("t01", "t02", "t03", "t04"),
		Type:      asttest.NewTypeId(asttest.NewIdent("T0"), classDeclT0.ToDeclarations()[0]),
	}
	varDeclT1s := &ast.VarDecl{
		IdentList: asttest.NewIdentList("t11", "t12", "t13", "t14"),
		Type:      asttest.NewTypeId(asttest.NewIdent("T1"), classDeclT1.ToDeclarations()[0]),
	}
	varDeclT2 := &ast.VarDecl{
		IdentList: asttest.NewIdentList("t22"),
		Type:      asttest.NewTypeId(asttest.NewIdent("T2"), classDeclT2.ToDeclarations()[0]),
	}
	varDeclT3 := &ast.VarDecl{
		IdentList: asttest.NewIdentList("t33"),
		Type:      asttest.NewTypeId(asttest.NewIdent("T3"), classDeclT3.ToDeclarations()[0]),
	}
	varDeclT4 := &ast.VarDecl{
		IdentList: asttest.NewIdentList("t44"),
		Type:      asttest.NewTypeId(asttest.NewIdent("T4"), classDeclT4.ToDeclarations()[0]),
	}

	writeln := func(arg string) *ast.Statement {
		writelnArg := asttest.NewExpression(asttest.NewString("'" + arg + "'"))
		return &ast.Statement{
			Body: &ast.CallStatement{
				Designator: asttest.NewDesignator(asttest.NewIdentRef("Writeln")),
				ExprList:   ast.ExprList{writelnArg},
			},
		}
	}

	createAndAssign := func(varName string, varDecl *ast.VarDecl, classDecl *ast.TypeDecl) *ast.Statement {
		return &ast.Statement{
			Body: &ast.AssignStatement{
				Designator: asttest.NewDesignator(
					asttest.NewQualId(
						varName,
						classDecl.ToDeclarations().Find(varName),
					),
				),
				Expression: ast.NewExpression(
					asttest.NewDesignatorFactor(
						asttest.NewDesignator(
							asttest.NewQualId(classDecl.Name, classDecl.ToDeclarations()[0]),
						).Add(
							asttest.NewDesignatorItemIdent("Create"),
						),
					),
				),
			},
		}
	}

	callFoo := func(varName string, varDecl *ast.VarDecl) *ast.Statement {
		return &ast.Statement{
			Body: &ast.CallStatement{
				Designator: asttest.NewDesignator(
					&ast.QualId{
						Ident: asttest.NewIdentRef(varName, varDecl.ToDeclarations()[0]),
					},
				).Add(
					asttest.NewDesignatorItemIdent(asttest.NewIdent("Foo")),
				),
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

	block := &ast.Block{
		DeclSections: ast.DeclSections{
			ast.VarSection{
				varDeclT0s,
				varDeclT1s,
				varDeclT2,
				varDeclT3,
				varDeclT4,
			},
		},
		Body: &ast.CompoundStmt{
			StmtList: ast.StmtList{
				writeln("T0"),
				createAndAssign("t01", varDeclT0s, classDeclT0),
				createAndAssign("t02", varDeclT0s, classDeclT0),
				createAndAssign("t03", varDeclT0s, classDeclT0),
				createAndAssign("t04", varDeclT0s, classDeclT0),
				callFoo("t01", varDeclT0s),
				callFoo("t02", varDeclT0s),
				callFoo("t03", varDeclT0s),
				callFoo("t04", varDeclT0s),
				writeln("T1"),
				createAndAssign("t11", varDeclT1s, classDeclT1),
				createAndAssign("t12", varDeclT1s, classDeclT1),
				createAndAssign("t13", varDeclT1s, classDeclT1),
				createAndAssign("t14", varDeclT1s, classDeclT1),
				callFoo("t11", varDeclT1s),
				callFoo("t12", varDeclT1s),
				callFoo("t13", varDeclT1s),
				callFoo("t14", varDeclT1s),
				writeln("T2"),
				createAndAssign("t22", varDeclT2, classDeclT2),
				callFoo("t22", varDeclT2),
				writeln("T3"),
				createAndAssign("t33", varDeclT3, classDeclT3),
				callFoo("t33", varDeclT3),
				writeln("T4"),
				createAndAssign("t44", varDeclT4, classDeclT4),
				callFoo("t44", varDeclT4),
				readln(),
			},
		},
	}

	expectedProject1 := &parser.Program{
		Program: &ast.Program{
			Path:  "Project1.dpr",
			Ident: asttest.NewIdent("Project1"),
			ProgramBlock: &ast.ProgramBlock{
				UsesClause: ast.UsesClause{
					asttest.NewUnitRef("SysUtils"),
					asttest.NewUnitRef("Unit1", "Unit1.pas"),
				},
				Block: block,
			},
		},
		Units: ast.Units{expectUnit1},
	}

	t.Run("Unit1", func(t *testing.T) {
		if !assert.Equal(t, expectUnit1, actualUnit1) {
			asttest.AssertUnit(t, expectUnit1, actualUnit1)
		}
	})
	t.Run("Project1", func(t *testing.T) {
		actualProject1.DeclMap = nil
		if !assert.Equal(t, expectedProject1, actualProject1) {
			asttest.AssertProgram(t, expectedProject1.Program, actualProject1.Program)
		}
	})
}
