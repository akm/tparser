package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestPointerType(t *testing.T) {

	RunBlockTest(t,
		"with pointer variable",
		[]rune(`
var
	X, Y: Integer;
	P: ^Integer;
begin
	X := 17;
	P := @X;
	Y := P^;
end;
`),
		&ast.Block{
			DeclSections: ast.DeclSections{
				ast.VarSection{
					{
						IdentList: asttest.NewIdentList("X", "Y"),
						Type:      asttest.NewOrdIdent("Integer"),
					},
					{
						IdentList: asttest.NewIdentList("P"),
						Type:      ast.NewCustomPointerType(asttest.NewOrdIdentWithIdent(asttest.NewIdent("Integer"))),
					},
				},
			},
			Body: &ast.CompoundStmt{
				StmtList: ast.StmtList{
					&ast.Statement{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdent("X")),
							Expression: asttest.NewExpression(asttest.NewNumber("17")),
						},
					},
					&ast.Statement{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdent("P")),
							Expression: asttest.NewExpression(&ast.Address{Designator: asttest.NewDesignator(asttest.NewIdent("X"))}),
						},
					},
					&ast.Statement{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdent("Y")),
							Expression: asttest.NewExpression(
								&ast.Designator{
									QualId: asttest.NewQualId(asttest.NewIdent("P")),
									Items:  ast.DesignatorItems{&ast.DesignatorItemDereference{}},
								},
							),
						},
					},
				},
			},
		},
	)

	RunBlockTest(t,
		"with pointer type declaration",
		[]rune(`
type
	PInteger = ^Integer;
var
	R: Single;	
	I: Integer;
	P: Pointer;
	PI: PInteger;
begin
	P := @R;
	PI := PInteger(P);
	I := PI^;
end;
`),
		&ast.Block{
			DeclSections: ast.DeclSections{
				ast.TypeSection{
					{
						Ident: asttest.NewIdent("PInteger"),
						Type:  ast.NewCustomPointerType(asttest.NewOrdIdentWithIdent(asttest.NewIdent("Integer"))),
					},
				},
				ast.VarSection{
					{
						IdentList: asttest.NewIdentList("R"),
						Type:      asttest.NewRealType("Single"),
					},
					{
						IdentList: asttest.NewIdentList("R"),
						Type:      asttest.NewOrdIdent("Single"),
					},
					{
						IdentList: asttest.NewIdentList("P"),
						Type:      ast.NewEmbeddedPointerType(asttest.NewIdent("Pointer")),
					},
					{
						IdentList: asttest.NewIdentList("P"),
						Type:      ast.NewTypeId(asttest.NewIdent("PInteger")),
					},
				},
			},
			Body: &ast.CompoundStmt{
				StmtList: ast.StmtList{
					&ast.Statement{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdent("P")),
							Expression: asttest.NewExpression(&ast.Address{Designator: asttest.NewDesignator(asttest.NewIdent("R"))}),
						},
					},
					&ast.Statement{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdent("PI")),
							Expression: asttest.NewExpression(
								&ast.TypeCast{
									TypeId: ast.NewTypeId(asttest.NewIdent("PInteger")),
									Expression: asttest.NewExpression(
										asttest.NewQualId(asttest.NewIdent("P")),
									),
								},
							),
						},
					},
					&ast.Statement{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdent("I")),
							Expression: asttest.NewExpression(
								&ast.Designator{
									QualId: asttest.NewQualId(asttest.NewIdent("PI")),
									Items:  ast.DesignatorItems{&ast.DesignatorItemDereference{}},
								},
							),
						},
					},
				},
			},
		},
	)

	// TODO tests for PChar, PAnsiChar, PWideChar
}
