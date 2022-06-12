package parsertest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestExportHeading(t *testing.T) {
	type pattern struct {
		text     string
		expected *ast.ExportedHeading
	}

	patterns := []pattern{}

	run := func(text string, expected *ast.ExportedHeading) {
		t.Run(text, func(t *testing.T) {
			patterns = append(patterns, pattern{text, expected})
			runes := []rune(text)
			parser := NewTestParser(&runes, parser.NewContext())
			parser.NextToken()
			res, err := parser.ParseExportedHeading()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		`PROCEDURE Proc0;`,
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("Proc0"),
			},
		},
	)

	run(
		`PROCEDURE Proc1(Param1: INTEGER);`,
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("Proc1"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("Param1", asttest.NewOrdIdent("INTEGER")),
				},
			},
		},
	)

	run(
		`procedure NumString(N: Integer; var S: string);`,
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("NumString"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("N", asttest.NewOrdIdent("Integer")),
					asttest.NewFormalParm("S", asttest.NewStringType("string"), &ast.FpoVar),
				},
			},
		},
	)

	run(
		"function WF: Integer;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:       ast.FtFunction,
				Ident:      asttest.NewIdent("WF"),
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"function Max(A: array of Real; N: Integer): Real;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("Max"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("A", asttest.NewArrayParameterType(asttest.NewRealType("Real"))),
					asttest.NewFormalParm("N", asttest.NewOrdIdent("Integer")),
				},
				ReturnType: asttest.NewRealType("Real"),
			},
		},
	)

	run(
		"function Power(X: Real; Y: Integer): Real;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("Power"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("X", asttest.NewRealType("Real")),
					asttest.NewFormalParm("Y", asttest.NewOrdIdent("Integer")),
				},
				ReturnType: asttest.NewRealType("Real"),
			},
		},
	)

	run(
		"function MyFunction(X, Y: Real): Real; cdecl;",
		&ast.ExportedHeading{
			Directives: []ast.Directive{ast.DrCdecl},
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("MyFunction"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"X", "Y"}, asttest.NewRealType("Real")),
				},
				ReturnType: asttest.NewRealType("Real"),
			},
		},
	)

	run(
		"function Calculate(X, Y: Integer): Real; forward;",
		&ast.ExportedHeading{
			Directives: []ast.Directive{ast.DrForward},
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("Calculate"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"X", "Y"}, asttest.NewOrdIdent("Integer")),
				},
				ReturnType: asttest.NewRealType("Real"),
			},
		},
	)

	run(
		"function printf(Format: PChar): Integer; cdecl; varargs;",
		&ast.ExportedHeading{
			Directives: []ast.Directive{ast.DrCdecl, ast.DrVarArgs},
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("printf"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("Format", "PChar"),
				},
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"procedure MoveWord(var Source, Dest; Count: Integer); external;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("MoveWord"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
					asttest.NewFormalParm("Count", asttest.NewOrdIdent("Integer")),
				},
			},
			Directives: []ast.Directive{ast.DrExternal},
		},
	)

	run(
		"function SomeFunction(S: string): string; external 'strlib.dll';",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("SomeFunction"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("S", asttest.NewStringType("string")),
				},
				ReturnType: asttest.NewStringType("string"),
			},
			Directives:      []ast.Directive{ast.DrExternal},
			ExternalOptions: &ast.ExternalOptions{LibraryName: "'strlib.dll'"},
		},
	)

	run(
		"function MessageBox(HWnd: Integer; Text, Caption: PChar; Flags: Integer): Integer; stdcall; external 'user32.dll' name 'MessageBoxA';",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("MessageBox"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm("HWnd", asttest.NewOrdIdent("Integer")),
					asttest.NewFormalParm([]string{"Text", "Caption"}, "PChar"),
					asttest.NewFormalParm("Flags", asttest.NewOrdIdent("Integer")),
				},
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
			Directives:      []ast.Directive{ast.DrStdcall, ast.DrExternal},
			ExternalOptions: &ast.ExternalOptions{LibraryName: "'user32.dll'", Name: ext.StringPtr("'MessageBoxA'")},
		},
	)

	run(
		"function Divide(X, Y: Real): Real; overload;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("Divide"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"X", "Y"}, asttest.NewRealType("Real")),
				},
				ReturnType: asttest.NewRealType("Real"),
			},
			Directives: []ast.Directive{ast.DrOverload},
		},
	)

	// These overloaded declarations are invalid because of they have same type and same length parameters
	//   function Cap(S: string): string; overload;
	//   procedure Cap(var Str: string); overload;
	run(
		"function DoubleByValue(X: Integer): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("DoubleByValue"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"X"}, asttest.NewOrdIdent("Integer")),
				},
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"function DoubleByRef(var X: Integer): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("DoubleByRef"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"X"}, asttest.NewOrdIdent("Integer"), &ast.FpoVar),
				},
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"function CompareStr(const S1, S2: string): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("CompareStr"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"S1", "S2"}, asttest.NewStringType("string"), &ast.FpoConst),
				},
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"procedure GetInfo(out Info: SomeRecordType);",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("GetInfo"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"Info"}, "SomeRecordType", &ast.FpoOut),
				},
			},
		},
	)

	run(
		"function Equal(var Source, Dest; Size: Integer): Boolean;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("Equal"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
					asttest.NewFormalParm([]string{"Size"}, asttest.NewOrdIdent("Integer")),
				},
				ReturnType: asttest.NewOrdIdent("Boolean"),
			},
		},
	)

	// {text: "procedure Check(S: string[20]);"} // syntax error
	run(
		"procedure Check(S: OpenString);",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("Check"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"S"}, "OpenString"),
				},
			},
		},
	)

	// {text: "procedure Sort(A: array[1..10] of Integer);"} // syntax error
	run(
		"function Find(A: array of Char): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("Find"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"A"}, asttest.NewArrayParameterType(asttest.NewOrdIdent("Char"))),
				},
				ReturnType: asttest.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"procedure Clear(var A: array of Real);",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("Clear"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"A"}, asttest.NewArrayParameterType(asttest.NewRealType("Real")), &ast.FpoVar),
				},
			},
		},
	)

	run(
		"function MakeStr(const Args: array of const): string;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("MakeStr"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"Args"}, asttest.NewArrayParameterType(nil), &ast.FpoConst),
				},
				ReturnType: asttest.NewStringType("string"),
			},
		},
	)

	run(
		"procedure FillArray(A: array of Integer; Value: Integer = 0);",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("FillArray"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm([]string{"A"}, asttest.NewArrayParameterType(asttest.NewOrdIdent("Integer"))),
					asttest.NewFormalParm(
						asttest.NewParameter([]string{"Value"}, asttest.NewOrdIdent("Integer"), asttest.NewNumber("0")),
					),
				},
			},
		},
	)

	run(
		"function MyFunction(X: Real = 3.5; Y: Real = 3.5): Real;",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: asttest.NewIdent("MyFunction"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm(
						asttest.NewParameter([]string{"X"}, asttest.NewRealType("Real"), asttest.NewNumber("3.5")),
					),
					asttest.NewFormalParm(
						asttest.NewParameter([]string{"Y"}, asttest.NewRealType("Real"), asttest.NewNumber("3.5")),
					),
				},
				ReturnType: asttest.NewRealType("Real"),
			},
		},
	)

	run(
		"procedure DoSomething(X: Real = 1.0; I: Integer = 0; S: string = '');",
		&ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: asttest.NewIdent("DoSomething"),
				FormalParameters: ast.FormalParameters{
					asttest.NewFormalParm(
						asttest.NewParameter([]string{"X"}, asttest.NewRealType("Real"), asttest.NewNumber("1.0")),
					),
					asttest.NewFormalParm(
						asttest.NewParameter([]string{"I"}, asttest.NewOrdIdent("Integer"), asttest.NewNumber("0")),
					),
					asttest.NewFormalParm(
						asttest.NewParameter([]string{"S"}, asttest.NewStringType("string"), asttest.NewString("''")),
					),
				},
			},
		},
	)
	// {text: "function MyFunction(X, Y: Real = 3.5): Real;"}, // syntax error
	// {text: "procedure MyProcedure(I: Integer = 1; S: string);"} // syntax error

	{
		headings := make([]string, len(patterns))
		decls := make([]ast.InterfaceDecl, len(patterns))
		for i, ptn := range patterns {
			headings[i] = ptn.text
			decls[i] = ptn.expected
		}

		RunUnitTest(t, "FunctionHeadings in unit",
			[]rune(fmt.Sprintf(`
UNIT Unit1;
INTERFACE
%s
IMPLEMENTATION
END.`, strings.Join(headings, "\n"))),
			&ast.Unit{
				Ident: asttest.NewIdent("Unit1"),
				InterfaceSection: &ast.InterfaceSection{
					InterfaceDecls: decls,
				},
				ImplementationSection: &ast.ImplementationSection{},
			},
		)

	}
}

func TestFormalParameters(t *testing.T) {
	run := func(text string, clearLocations bool, expected ast.FormalParameters) {
		t.Run(text, func(t *testing.T) {
			runes := []rune(text)
			parser := NewTestParser(&runes, parser.NewContext())
			parser.NextToken()
			res, err := parser.ParseFormalParameters()
			if assert.NoError(t, err) {
				if clearLocations {
					asttest.ClearLocations(t, res)
				}
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"(X, Y: Real)", true,
		ast.FormalParameters{
			asttest.NewFormalParm([]string{"X", "Y"}, asttest.NewRealType("Real")),
		},
	)
	run(
		"(var S: string; X: Integer)", true,
		ast.FormalParameters{
			asttest.NewFormalParm("S", asttest.NewStringType("string"), &ast.FpoVar),
			asttest.NewFormalParm("X", asttest.NewOrdIdent("Integer")),
		},
	)
	run(
		"(HWnd: Integer; Text, Caption: PChar; PChar: Integer)", false,
		ast.FormalParameters{
			asttest.NewFormalParm(
				asttest.NewIdent("HWnd", asttest.NewIdentLocation(1, 2, 1, 6)),
				asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(1, 8, 7, 15))),
			),
			asttest.NewFormalParm(
				[]*ast.Ident{
					asttest.NewIdent("Text", asttest.NewIdentLocation(1, 17, 16, 21)),
					asttest.NewIdent("Caption", asttest.NewIdentLocation(1, 23, 22, 30)),
				},
				asttest.NewIdent("PChar", asttest.NewIdentLocation(1, 32, 31, 37)),
			),
			asttest.NewFormalParm(
				asttest.NewIdent("PChar", asttest.NewIdentLocation(1, 39, 38, 44)),
				asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(1, 46, 45, 53))),
			),
		},
	)
	run(
		"(const P; I: Integer)", true,
		ast.FormalParameters{
			asttest.NewFormalParm([]string{"P"}, nil, &ast.FpoConst),
			asttest.NewFormalParm([]string{"I"}, asttest.NewOrdIdent("Integer")),
		},
	)
}
