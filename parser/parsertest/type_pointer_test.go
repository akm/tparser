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
		func() *ast.Block {

			varDeclXY := &ast.VarDecl{
				IdentList: asttest.NewIdentList("X", "Y"),
				Type:      asttest.NewOrdIdent("Integer"),
			}
			varDeclP := &ast.VarDecl{
				IdentList: asttest.NewIdentList("P"),
				Type:      ast.NewCustomPointerType(asttest.NewOrdIdentWithIdent(asttest.NewIdent("Integer"))),
			}
			declX := varDeclXY.ToDeclarations()[0]
			declY := varDeclXY.ToDeclarations()[1]
			declP := varDeclP.ToDeclarations()[0]
			return &ast.Block{
				DeclSections: ast.DeclSections{
					ast.VarSection{varDeclXY, varDeclP},
				},
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("X", declX)),
								Expression: asttest.NewExpression(asttest.NewNumber("17")),
							},
						},
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("P", declP)),
								Expression: asttest.NewExpression(&ast.Address{Designator: asttest.NewDesignator(asttest.NewIdentRef("X", declX))}),
							},
						},
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("Y", declY)),
								Expression: asttest.NewExpression(
									&ast.Designator{
										QualId: asttest.NewQualId(asttest.NewIdentRef("P", declP)),
										Items:  ast.DesignatorItems{&ast.DesignatorItemDereference{}},
									},
								),
							},
						},
					},
				},
			}
		}(),
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
		func() *ast.Block {
			typeDeclPInteger := &ast.TypeDecl{
				Ident: asttest.NewIdent("PInteger"),
				Type:  ast.NewCustomPointerType(asttest.NewOrdIdentWithIdent(asttest.NewIdent("Integer"))),
			}
			varDeclR := &ast.VarDecl{
				IdentList: asttest.NewIdentList("R"),
				Type:      asttest.NewRealType("Single"),
			}
			varDeclI := &ast.VarDecl{
				IdentList: asttest.NewIdentList("I"),
				Type:      asttest.NewOrdIdent("Integer"),
			}
			varDeclP := &ast.VarDecl{
				IdentList: asttest.NewIdentList("P"),
				Type:      ast.NewEmbeddedPointerType(asttest.NewIdent("Pointer")),
			}
			varDeclPI := &ast.VarDecl{
				IdentList: asttest.NewIdentList("PI"),
				Type:      ast.NewTypeId(asttest.NewIdent("PInteger")),
			}
			declPInteger := typeDeclPInteger.ToDeclarations()[0]
			declR := varDeclR.ToDeclarations()[0]
			declI := varDeclI.ToDeclarations()[0]
			declP := varDeclP.ToDeclarations()[0]
			declPI := varDeclPI.ToDeclarations()[0]
			return &ast.Block{
				DeclSections: ast.DeclSections{
					ast.TypeSection{typeDeclPInteger},
					ast.VarSection{varDeclR, varDeclI, varDeclP, varDeclPI},
				},
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("P", declP)),
								Expression: asttest.NewExpression(&ast.Address{Designator: asttest.NewDesignator(asttest.NewIdentRef("R", declR))}),
							},
						},
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("PI", declPI)),
								Expression: asttest.NewExpression(
									&ast.TypeCast{
										TypeId: ast.NewTypeId(asttest.NewIdentRef("PInteger", declPInteger)),
										Expression: asttest.NewExpression(
											asttest.NewQualId(asttest.NewIdentRef("P", declP)),
										),
									},
								),
							},
						},
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("I", declI)),
								Expression: asttest.NewExpression(
									&ast.Designator{
										QualId: asttest.NewQualId(asttest.NewIdentRef("PI", declPI)),
										Items:  ast.DesignatorItems{&ast.DesignatorItemDereference{}},
									},
								),
							},
						},
					},
				},
			}
		}(),
	)

	// TODO tests for PChar, PAnsiChar, PWideChar
}
