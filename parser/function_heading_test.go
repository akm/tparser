package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ext"
	"github.com/stretchr/testify/assert"
)

func TestExportHeading(t *testing.T) {
	run := func(text string, expected *ast.ExportedHeading) {
		t.Run(text, func(t *testing.T) {
			runes := []rune(text)
			parser := NewParser(&runes, NewContext())
			parser.NextToken()
			res, err := parser.ParseExportedHeading()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	type pattern struct {
		text     string
		expected *ast.ExportedHeading
	}

	patterns := []pattern{
		{
			text: `PROCEDURE Proc0;`,
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("Proc0"),
				},
			},
		},
		{
			text: `PROCEDURE Proc1(Param1: INTEGER);`,
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("Proc1"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("Param1", ast.NewOrdIdent("INTEGER")),
					},
				},
			},
		},
		{
			text: `procedure NumString(N: Integer; var S: string);`,
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("NumString"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("N", ast.NewOrdIdent("Integer")),
						ast.NewFormalParm("S", ast.NewStringType("STRING"), &ast.FpoVar),
					},
				},
			},
		},
		{
			text: "function WF: Integer;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:       ast.FtFunction,
					Ident:      ast.Ident("WF"),
					ReturnType: ast.NewOrdIdent("Integer"),
				},
			},
		},
		{
			text: "function Max(A: array of Real; N: Integer): Real;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("Max"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("A", ast.NewArrayParameterType(ast.NewRealType("Real"))),
						ast.NewFormalParm("N", ast.NewOrdIdent("Integer")),
					},
					ReturnType: ast.NewRealType("Real"),
				},
			},
		},
		{
			text: "function Power(X: Real; Y: Integer): Real;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("Power"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("X", ast.NewRealType("Real")),
						ast.NewFormalParm("Y", ast.NewOrdIdent("Integer")),
					},
					ReturnType: ast.NewRealType("Real"),
				},
			},
		},
		{
			text: "function MyFunction(X, Y: Real): Real; cdecl;",
			expected: &ast.ExportedHeading{
				Directive: []ast.Directive{ast.DrCdecl},
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("MyFunction"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"X", "Y"}, ast.NewRealType("Real")),
					},
					ReturnType: ast.NewRealType("Real"),
				},
			},
		},
		{
			text: "function Calculate(X, Y: Integer): Real; forward;",
			expected: &ast.ExportedHeading{
				Directive: []ast.Directive{ast.DrForward},
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("Calculate"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"X", "Y"}, ast.NewOrdIdent("Integer")),
					},
					ReturnType: ast.NewRealType("Real"),
				},
			},
		},
		{
			text: "function printf(Format: PChar): Integer; cdecl; varargs;",
			expected: &ast.ExportedHeading{
				Directive: []ast.Directive{ast.DrCdecl, ast.DrVarArgs},
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("printf"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("Format", "PChar"),
					},
					ReturnType: ast.NewOrdIdent("Integer"),
				},
			},
		},
		{
			text: "procedure MoveWord(var Source, Dest; Count: Integer); external;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("MoveWord"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
						ast.NewFormalParm("Count", ast.NewOrdIdent("Integer"), &ast.FpoVar),
					},
				},
				Directive: []ast.Directive{ast.DrExternal},
			},
		},
		{
			text: "function SomeFunction(S: string): string; external 'strlib.dll';",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("SomeFunction"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("S", ast.NewStringType("STRING")),
					},
					ReturnType: ast.NewStringType("STRING"),
				},
				Directive:       []ast.Directive{ast.DrExternal},
				ExternalOptions: &ast.ExternalOptions{LibraryName: "'strlib.dll'"},
			},
		},
		{
			text: "function MessageBox(HWnd: Integer; Text, Caption: PChar; Flags: Integer): Integer; stdcall; external 'user32.dll' name 'MessageBoxA';",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("MessageBox"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm("HWnd", ast.NewOrdIdent("Integer")),
						ast.NewFormalParm([]string{"Text", "Caption"}, "PChar"),
						ast.NewFormalParm("PChar", ast.NewOrdIdent("Integer")),
					},
					ReturnType: ast.NewOrdIdent("Integer"),
				},
				Directive:       []ast.Directive{ast.DrStdcall, ast.DrExternal},
				ExternalOptions: &ast.ExternalOptions{LibraryName: "'user32.dll'", Name: ext.StringPtr("'MessageBoxA'")},
			},
		},
		{
			text: "function Divide(X, Y: Real): Real; overload;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("Divide"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"X", "Y"}, ast.NewRealType("Real")),
					},
					ReturnType: ast.NewRealType("Real"),
				},
				Directive: []ast.Directive{ast.DrOverload},
			},
		},
		// These overloaded declarations are invalid because of they have same type and same length parameters
		//   function Cap(S: string): string; overload;
		//   procedure Cap(var Str: string); overload;
		{
			text: "function DoubleByValue(X: Integer): Integer;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("DoubleByValue"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"X"}, ast.NewOrdIdent("Integer")),
					},
					ReturnType: ast.NewOrdIdent("Integer"),
				},
			},
		},
		{
			text: "function DoubleByRef(var X: Integer): Integer;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("DoubleByRef"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"X"}, ast.NewOrdIdent("Integer"), &ast.FpoVar),
					},
					ReturnType: ast.NewOrdIdent("Integer"),
				},
			},
		},
		{
			text: "function CompareStr(const S1, S2: string): Integer;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("CompareStr"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"S1", "S2"}, ast.NewStringType("STRING"), &ast.FpoConst),
					},
					ReturnType: ast.NewOrdIdent("Integer"),
				},
			},
		},
		{
			text: "procedure GetInfo(out Info: SomeRecordType);",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("GetInfo"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"Info"}, "SomeRecordType", &ast.FpoOut),
					},
				},
			},
		},
		{
			text: "function Equal(var Source, Dest; Size: Integer): Boolean;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("Equal"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
						ast.NewFormalParm([]string{"Size"}, ast.NewOrdIdent("Integer")),
					},
					ReturnType: ast.NewOrdIdent("Boolean"),
				},
			},
		},
		// {text: "procedure Check(S: string[20]);"} // syntax error
		{
			text: "procedure Check(S: OpenString);",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("Check"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"S"}, "OpenString"),
					},
				},
			},
		},
		// {text: "procedure Sort(A: array[1..10] of Integer);"} // syntax error
		{
			text: "function Find(A: array of Char): Integer;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("Find"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"A"}, ast.NewArrayParameterType(ast.NewOrdIdent("Char"))),
					},
					ReturnType: ast.NewOrdIdent("Integer"),
				},
			},
		},
		{
			text: "procedure Clear(var A: array of Real);",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("Clear"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"A"}, ast.NewArrayParameterType(ast.NewRealType("Real")), &ast.FpoVar),
					},
				},
			},
		},
		{
			text: "function MakeStr(const Args: array of const): string;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("MakeStr"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"Args"}, nil, &ast.FpoConst),
					},
					ReturnType: ast.NewStringType("string"),
				},
			},
		},
		{
			text: "procedure FillArray(A: array of Integer; Value: Integer = 0);",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("FillArray"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"A"}, ast.NewArrayParameterType(ast.NewOrdIdent("Integer"))),
						ast.NewFormalParm(
							ast.NewParameter([]string{"Value"}, ast.NewOrdIdent("Integer"), ast.NewNumber("0")),
						),
					},
				},
			},
		},
		{
			text: "function MyFunction(X: Real = 3.5; Y: Real = 3.5): Real;",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: ast.Ident("MyFunction"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm(
							ast.NewParameter([]string{"X"}, ast.NewRealType("Real"), ast.NewNumber("3.5")),
						),
						ast.NewFormalParm(
							ast.NewParameter([]string{"Y"}, ast.NewRealType("Real"), ast.NewNumber("3.5")),
						),
					},
					ReturnType: ast.NewRealType("Real"),
				},
			},
		},
		{
			text: "procedure DoSomething(X: Real = 1.0; I: Integer = 0; S: string = '');",
			expected: &ast.ExportedHeading{
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("DoSomething"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm(
							ast.NewParameter([]string{"X"}, ast.NewRealType("Real"), ast.NewNumber("1.0")),
						),
						ast.NewFormalParm(
							ast.NewParameter([]string{"I"}, ast.NewOrdIdent("Integer"), ast.NewNumber("0")),
						),
						ast.NewFormalParm(
							ast.NewParameter([]string{"S"}, ast.NewStringType("STRING"), ast.NewString("''")),
						),
					},
				},
			},
		},
		// {text: "function MyFunction(X, Y: Real = 3.5): Real;"}, // syntax error
		// {text: "procedure MyProcedure(I: Integer = 1; S: string);"} // syntax error

	}

	for _, ptn := range patterns {
		run(ptn.text, ptn.expected)
	}

	t.Run("FunctionHeadings in unit", func(t *testing.T) {
		headings := make([]string, len(patterns))
		decls := make([]ast.InterfaceDecl, len(patterns))
		for i, ptn := range patterns {
			headings[i] = ptn.text
			decls[i] = ptn.expected
		}

		unitText := []rune(fmt.Sprintf(`UNIT Unit1;
		INTERFACE
		%s
		IMPLEMENTATION
		END.`, strings.Join(headings, "\n")))

		parser := NewParser(&unitText)
		parser.NextToken()
		res, err := parser.ParseUnit()
		if assert.NoError(t, err) {
			assert.Equal(t,
				&ast.Unit{
					Ident: ast.Ident("Unit1"),
					InterfaceSection: &ast.InterfaceSection{
						InterfaceDecls: decls,
					},
					ImplementationSection: &ast.ImplementationSection{},
				},
				res,
			)
		}
	})
}

func TestFormalParameters(t *testing.T) {
	run := func(text string, expected ast.FormalParameters) {
		t.Run(text, func(t *testing.T) {
			runes := []rune(text)
			parser := NewParser(&runes, NewContext())
			parser.NextToken()
			res, err := parser.ParseFormalParameters()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"(X, Y: Real)",
		ast.FormalParameters{
			ast.NewFormalParm([]string{"X", "Y"}, ast.NewRealType("Real")),
		},
	)
	run(
		"(var S: string; X: Integer)",
		ast.FormalParameters{
			ast.NewFormalParm("S", ast.NewStringType("STRING"), &ast.FpoVar),
			ast.NewFormalParm("X", ast.NewOrdIdent("Integer")),
		},
	)
	run(
		"(HWnd: Integer; Text, Caption: PChar; PChar: Integer)",
		ast.FormalParameters{
			ast.NewFormalParm("HWnd", ast.NewOrdIdent("Integer")),
			ast.NewFormalParm([]string{"Text", "Caption"}, "PChar"),
			ast.NewFormalParm("PChar", ast.NewOrdIdent("Integer")),
		},
	)
	run(
		"(const P; I: Integer)",
		ast.FormalParameters{
			ast.NewFormalParm([]string{"P"}, nil, &ast.FpoConst),
			ast.NewFormalParm([]string{"I"}, ast.NewOrdIdent("Integer")),
		},
	)
}
