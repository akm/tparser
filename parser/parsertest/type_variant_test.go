package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestVariantType(t *testing.T) {

	RunBlockTest(t,
		"variant variables",
		[]rune(`
var
	V1, V2, V3, V4, V5: Variant;
	I: Integer;
	D: Double;
	S: string;
begin
	V1 := 1; { integer value }
	V2 := 1234.5678; { real value }
	V3 := 'Hello world!'; { string value }
	V4 := '1000'; { string value }
	V5 := V1 + V2 + V4; { real value 2235.5678}
	I := V1; { I = 1 (integer value) }
	D := V2; { D = 1234.5678 (real value) }
	S := V3; { S = 'Hello world!' (string value) }
	I := V4; { I = 1000 (integer value) }
	S := V5; { S = '2235.5678' (string value) }
end;
`),
		func() *ast.Block {
			varDeclV1ToV5 := &ast.VarDecl{
				IdentList: asttest.NewIdentList("V1", "V2", "V3", "V4", "V5"),
				Type:      asttest.NewVariantType("Variant"),
			}
			varDeclI := &ast.VarDecl{
				IdentList: asttest.NewIdentList("I"),
				Type:      asttest.NewOrdIdent("Integer"),
			}
			varDeclD := &ast.VarDecl{
				IdentList: asttest.NewIdentList("D"),
				Type:      asttest.NewRealType("Double"),
			}
			varDeclS := &ast.VarDecl{
				IdentList: asttest.NewIdentList("S"),
				Type:      asttest.NewStringType("string"),
			}
			declVs := varDeclV1ToV5.ToDeclarations()
			declV1 := declVs[0]
			declV2 := declVs[1]
			declV3 := declVs[2]
			declV4 := declVs[3]
			declV5 := declVs[4]
			declI := varDeclI.ToDeclarations()[0]
			declD := varDeclD.ToDeclarations()[0]
			declS := varDeclS.ToDeclarations()[0]
			return &ast.Block{
				DeclSections: ast.DeclSections{
					ast.VarSection{varDeclV1ToV5, varDeclI, varDeclD, varDeclS},
				},
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("V1", declV1)),
							Expression: asttest.NewExpression(asttest.NewNumber("1")),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("V2", declV2)),
							Expression: asttest.NewExpression(asttest.NewNumber("1234.5678")),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("V3", declV3)),
							Expression: asttest.NewExpression(asttest.NewString("'Hello world!'")),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("V4", declV4)),
							Expression: asttest.NewExpression(asttest.NewString("'1000'")),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("V5", declV5)),
							Expression: asttest.NewExpression(
								&ast.SimpleExpression{
									Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewQualId("V1", declV1))},
									AddOpTerms: []*ast.AddOpTerm{
										{AddOp: "+", Term: asttest.NewTerm(asttest.NewQualId("V2", declV2))},
										{AddOp: "+", Term: asttest.NewTerm(asttest.NewQualId("V4", declV4))},
									},
								},
							),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("I", declI)),
							Expression: asttest.NewExpression(asttest.NewQualId("V1", declV1)),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("D", declD)),
							Expression: asttest.NewExpression(asttest.NewQualId("V2", declV2)),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("S", declS)),
							Expression: asttest.NewExpression(asttest.NewQualId("V3", declV3)),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("I", declI)),
							Expression: asttest.NewExpression(asttest.NewQualId("V4", declV4)),
						}},
						&ast.Statement{Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator(asttest.NewIdentRef("S", declS)),
							Expression: asttest.NewExpression(asttest.NewQualId("V5", declV5)),
						}},
					},
				},
			}
		}(),
	)

}
