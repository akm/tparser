package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestProcedureType(t *testing.T) {

	RunVarSectionTest(t,
		"var with procedure type definition",
		[]rune(`var F: function(X,Y: Integer): Integer;`),
		ast.VarSection{
			&ast.VarDecl{
				IdentList: asttest.NewIdentList("F"),
				Type: &ast.ProcedureType{
					FunctionType: ast.FtFunction,
					FormalParameters: ast.FormalParameters{
						asttest.NewFormalParm([]string{"X", "Y"}, asttest.NewOrdIdent("Integer")),
					},
					ReturnType: asttest.NewOrdIdent("Integer"),
				},
			},
		},
	)

	RunUnitTest(t,
		"procedure type declaration and vars",
		[]rune(`UNIT U1;
interface
type
	TIntegerFunction = function: Integer;
	TProcedure = procedure;
	TStrProc = procedure(const S: string);
	TMathFunc = function(X: Double): Double;
var
	F: TIntegerFunction;
	Proc: TProcedure;
	SP: TStrProc;
	M: TMathFunc;

procedure FuncProc(P: TIntegerFunction);	

implementation
end.`),
		func() *ast.Unit {
			typeDeclTIntegerFunction := &ast.TypeDecl{
				Ident: asttest.NewIdent("TIntegerFunction"),
				Type: &ast.ProcedureType{
					FunctionType: ast.FtFunction,
					ReturnType:   asttest.NewOrdIdent("Integer"),
				},
			}
			typeDeclTProcedure := &ast.TypeDecl{
				Ident: asttest.NewIdent("TProcedure"),
				Type:  &ast.ProcedureType{FunctionType: ast.FtProcedure},
			}
			typeDeclTStrProc := &ast.TypeDecl{
				Ident: asttest.NewIdent("TStrProc"),
				Type: &ast.ProcedureType{
					FunctionType: ast.FtProcedure,
					FormalParameters: ast.FormalParameters{
						asttest.NewFormalParm("S", asttest.NewStringType("string"), "CONST"),
					},
				},
			}
			typeDeclTMathFunc := &ast.TypeDecl{
				Ident: asttest.NewIdent("TMathFunc"),
				Type: &ast.ProcedureType{
					FunctionType: ast.FtFunction,
					FormalParameters: ast.FormalParameters{
						asttest.NewFormalParm("X", asttest.NewRealType("Double")),
					},
					ReturnType: asttest.NewRealType("Double"),
				},
			}

			declTIntegerFunction := typeDeclTIntegerFunction.ToDeclarations()[0]
			declTProcedure := typeDeclTProcedure.ToDeclarations()[0]
			declTStrProc := typeDeclTStrProc.ToDeclarations()[0]
			declTMathFunc := typeDeclTMathFunc.ToDeclarations()[0]

			return &ast.Unit{
				Ident: asttest.NewIdent("U1"),
				InterfaceSection: &ast.InterfaceSection{
					InterfaceDecls: ast.InterfaceDecls{
						ast.TypeSection{
							typeDeclTIntegerFunction,
							typeDeclTProcedure,
							typeDeclTStrProc,
							typeDeclTMathFunc,
						},
						ast.VarSection{
							&ast.VarDecl{IdentList: asttest.NewIdentList("F"), Type: asttest.NewTypeId("TIntegerFunction", declTIntegerFunction)},
							&ast.VarDecl{IdentList: asttest.NewIdentList("Proc"), Type: asttest.NewTypeId("TProcedure", declTProcedure)},
							&ast.VarDecl{IdentList: asttest.NewIdentList("SP"), Type: asttest.NewTypeId("TStrProc", declTStrProc)},
							&ast.VarDecl{IdentList: asttest.NewIdentList("M"), Type: asttest.NewTypeId("TMathFunc", declTMathFunc)},
						},
						&ast.ExportedHeading{
							FunctionHeading: &ast.FunctionHeading{
								Type:  ast.FtProcedure,
								Ident: asttest.NewIdent("FuncProc"),
								FormalParameters: ast.FormalParameters{
									asttest.NewFormalParm("P", asttest.NewTypeId("TIntegerFunction", declTIntegerFunction)),
								},
							},
						},
					},
				},
				ImplementationSection: &ast.ImplementationSection{},
			}
		}(),
	)

	RunUnitTest(t,
		"procedures with same name parameter",
		[]rune(`UNIT U1;
interface
type
	TStrProc1 = procedure(S: string);
	TStrProc2 = procedure(S: string);

procedure Proc1(S: string);
procedure Proc2(S: string);

implementation

procedure Proc1(S: string);
begin
	Writeln(S);
end;

procedure Proc2(S: string);
begin
	Writeln(S);
end;

end.`),
		func() *ast.Unit {
			typeDeclTStrProc1 := &ast.TypeDecl{
				Ident: asttest.NewIdent("TStrProc1"),
				Type: &ast.ProcedureType{
					FunctionType:     ast.FtProcedure,
					FormalParameters: ast.FormalParameters{asttest.NewFormalParm("S", asttest.NewStringType("string"))},
				},
			}
			typeDeclTStrProc2 := &ast.TypeDecl{
				Ident: asttest.NewIdent("TStrProc2"),
				Type: &ast.ProcedureType{
					FunctionType:     ast.FtProcedure,
					FormalParameters: ast.FormalParameters{asttest.NewFormalParm("S", asttest.NewStringType("string"))},
				},
			}

			paramSinProc1 := asttest.NewFormalParm("S", asttest.NewStringType("string"))
			paramSinProc2 := asttest.NewFormalParm("S", asttest.NewStringType("string"))

			return &ast.Unit{
				Ident: asttest.NewIdent("U1"),
				InterfaceSection: &ast.InterfaceSection{
					InterfaceDecls: ast.InterfaceDecls{
						ast.TypeSection{
							typeDeclTStrProc1,
							typeDeclTStrProc2,
						},
						&ast.ExportedHeading{
							FunctionHeading: &ast.FunctionHeading{
								Type:             ast.FtProcedure,
								Ident:            asttest.NewIdent("Proc1"),
								FormalParameters: ast.FormalParameters{asttest.NewFormalParm("S", asttest.NewStringType("string"))},
							},
						},
						&ast.ExportedHeading{
							FunctionHeading: &ast.FunctionHeading{
								Type:             ast.FtProcedure,
								Ident:            asttest.NewIdent("Proc2"),
								FormalParameters: ast.FormalParameters{asttest.NewFormalParm("S", asttest.NewStringType("string"))},
							},
						},
					},
				},
				ImplementationSection: &ast.ImplementationSection{
					DeclSections: ast.DeclSections{
						&ast.FunctionDecl{
							FunctionHeading: &ast.FunctionHeading{
								Type:             ast.FtProcedure,
								Ident:            asttest.NewIdent("Proc1"),
								FormalParameters: ast.FormalParameters{paramSinProc1},
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
														&ast.DesignatorFactor{
															Designator: asttest.NewDesignator(
																&ast.QualId{Ident: asttest.NewIdentRef("S", paramSinProc1.ToDeclarations()[0])},
															),
														},
													),
												},
											},
										},
									},
								},
							},
						},
						&ast.FunctionDecl{
							FunctionHeading: &ast.FunctionHeading{
								Type:             ast.FtProcedure,
								Ident:            asttest.NewIdent("Proc2"),
								FormalParameters: ast.FormalParameters{asttest.NewFormalParm("S", asttest.NewStringType("string"))},
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
														&ast.DesignatorFactor{
															Designator: asttest.NewDesignator(
																&ast.QualId{Ident: asttest.NewIdentRef("S", paramSinProc2.ToDeclarations()[0])},
															),
														},
													),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
		}(),
	)

	RunTypeSection(t,
		"Object method type declaration",
		[]rune(`
type
	TMethod = procedure of object;
	TNotifyEvent = procedure(Sender: TObject) of object;
`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: asttest.NewIdent("TMethod"),
				Type: &ast.ProcedureType{
					FunctionType: ast.FtProcedure,
					OfObject:     true,
				},
			},
			&ast.TypeDecl{
				Ident: asttest.NewIdent("TNotifyEvent"),
				Type: &ast.ProcedureType{
					FunctionType: ast.FtProcedure,
					FormalParameters: ast.FormalParameters{
						asttest.NewFormalParm("Sender", asttest.NewTypeId("TObject")),
					},
					OfObject: true,
				},
			},
		},
	)

}
