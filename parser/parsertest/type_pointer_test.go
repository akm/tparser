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
			declPInteger := typeDeclPInteger.ToDeclarations()[0]

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
				Type:      ast.NewTypeId(asttest.NewIdent("PInteger"), declPInteger),
			}
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
										TypeId: ast.NewTypeId(asttest.NewIdent("PInteger"), declPInteger),
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

	RunBlockTest(t,
		"with PChar",
		[]rune(`
var
	MyArray: array[0..32] of Char;
	MyPointer: PChar;
begin
	MyArray := 'Hello';
	MyPointer := MyArray;
	SomeProcedure(MyArray);
	SomeProcedure(MyPointer);
end;
`),
		func() *ast.Block {
			varDeclMyArray := &ast.VarDecl{
				IdentList: asttest.NewIdentList("MyArray"),
				Type: &ast.ArrayType{
					IndexTypes: []ast.OrdinalType{
						&ast.SubrangeType{
							Low:  asttest.NewConstExpr(asttest.NewNumber("0")),
							High: asttest.NewConstExpr(asttest.NewNumber("32")),
						},
					},
					BaseType: ast.NewOrdIdent(asttest.NewIdent("Char")),
				},
			}
			varDeclMyPointer := &ast.VarDecl{
				IdentList: asttest.NewIdentList("MyPointer"),
				Type:      ast.NewEmbeddedPointerType(asttest.NewIdent("PChar")),
			}
			declMyArray := varDeclMyArray.ToDeclarations()[0]
			declMyPointer := varDeclMyPointer.ToDeclarations()[0]
			return &ast.Block{
				DeclSections: ast.DeclSections{
					ast.VarSection{varDeclMyArray, varDeclMyPointer},
				},
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("MyArray", declMyArray)),
								Expression: asttest.NewExpression(asttest.NewString("'Hello'")),
							},
						},
						&ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("MyPointer", declMyPointer)),
								Expression: asttest.NewExpression(asttest.NewIdentRef("MyArray", declMyArray)),
							},
						},
						&ast.Statement{
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("SomeProcedure")),
								ExprList:   ast.ExprList{asttest.NewExpression(asttest.NewIdentRef("MyArray", declMyArray))},
							},
						},
						&ast.Statement{
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator(asttest.NewIdentRef("SomeProcedure")),
								ExprList:   ast.ExprList{asttest.NewExpression(asttest.NewIdentRef("MyPointer", declMyPointer))},
							},
						},
					},
				},
			}
		}(),
	)

	NewTypeSectionTestRunner(t,
		"Forward declaration of PointerType",
		[]rune(`
type
	PGUID = ^TGUID;
	TGUID = packed record
		D1: Longword;
		D2: Word;
		D3: Word;
		D4: array[0..7] of Byte;
	end;
`),
		func() ast.TypeSection {
			tguidType := &ast.RecType{
				Packed: true,
				FieldList: &ast.FieldList{
					FieldDecls: ast.FieldDecls{
						{
							IdentList: asttest.NewIdentList("D1"),
							Type:      asttest.NewOrdIdent("Longword"),
						},
						{
							IdentList: asttest.NewIdentList("D2"),
							Type:      asttest.NewOrdIdent("Word"),
						},
						{
							IdentList: asttest.NewIdentList("D3"),
							Type:      asttest.NewOrdIdent("Word"),
						},
						{
							IdentList: asttest.NewIdentList("D4"),
							Type: &ast.ArrayType{
								IndexTypes: []ast.OrdinalType{
									&ast.SubrangeType{
										Low:  asttest.NewConstExpr(asttest.NewNumber("0")),
										High: asttest.NewConstExpr(asttest.NewNumber("7")),
									},
								},
								BaseType: ast.NewOrdIdent(asttest.NewIdent("Byte")),
							},
						},
					},
				},
			}

			recDecl := &ast.TypeDecl{
				Ident: asttest.NewIdent("TGUID"),
				Type:  tguidType,
			}
			pointerDecl := &ast.TypeDecl{
				Ident: asttest.NewIdent("PGUID"),
				Type:  ast.NewCustomPointerType(asttest.NewTypeId("TGUID", recDecl.ToDeclarations()[0])),
			}

			return ast.TypeSection{
				pointerDecl,
				recDecl,
			}
		}(),
	).Run()

}
