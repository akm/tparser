package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akm/tparser/ast"
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
						ast.NewFormalParm("Param1", "Integer"),
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
						ast.NewFormalParm("N", "Integer"),
						ast.NewFormalParm("S", "string", &ast.FpoVar),
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
					ReturnType: &ast.TypeId{Ident: ast.Ident("INTEGER")},
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
						ast.NewFormalParm("A", ast.NewArrayParameterType("Real")),
					},
					ReturnType: &ast.TypeId{Ident: ast.Ident("Real")},
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
						ast.NewFormalParm("X", "Real"),
						ast.NewFormalParm("Y", "Integer"),
					},
					ReturnType: &ast.TypeId{Ident: ast.Ident("Real")},
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
						ast.NewFormalParm([]string{"X", "Y"}, "Real"),
					},
					ReturnType: &ast.TypeId{Ident: ast.Ident("Real")},
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
						ast.NewFormalParm([]string{"X", "Y"}, "Integer"),
					},
					ReturnType: &ast.TypeId{Ident: ast.Ident("Real")},
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
					ReturnType: &ast.TypeId{Ident: ast.Ident("Integer")},
				},
			},
		},
		{
			text: "procedure MoveWord(var Source, Dest; Count: Integer); external;",
			expected: &ast.ExportedHeading{
				Directive: []ast.Directive{ast.DrExternal},
				FunctionHeading: ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: ast.Ident("MoveWord"),
					FormalParameters: ast.FormalParameters{
						ast.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
						ast.NewFormalParm("Count", "Integer", &ast.FpoVar),
					},
				},
			},
		},
		{text: "function SomeFunction(S: string): string; external 'strlib.dll';"},
		{text: "function MessageBox(HWnd: Integer; Text, Caption: PChar; Flags: Integer): Integer; stdcall; external 'user32.dll' name 'MessageBoxA';"},
		{text: "function Divide(X, Y: Real): Real; overload;"},
		// These overloaded declarations are invalid because of they have same type and same length parameters
		//   function Cap(S: string): string; overload;
		//   procedure Cap(var Str: string); overload;
		{text: "function DoubleByValue(X: Integer): Integer;"},
		{text: "function DoubleByRef(var X: Integer): Integer;"},
		{text: "function CompareStr(const S1, S2: string): Integer;"},
		{text: "procedure GetInfo(out Info: SomeRecordType);"},
		{text: "function Equal(var Source, Dest; Size: Integer): Boolean;"},
		// {text: "procedure Check(S: string[20]);"} // syntax error
		{text: "procedure Check(S: OpenString);"},
		// {text: "procedure Sort(A: array[1..10] of Integer);"} // syntax error
		{text: "function Find(A: array of Char): Integer;"},
		{text: "procedure Clear(var A: array of Real);"},
		{text: "function MakeStr(const Args: array of const): string;"},
		{text: "procedure FillArray(A: array of Integer; Value: Integer = 0);"},
		{text: "function MyFunction(X: Real = 3.5; Y: Real = 3.5): Real;"},
		{text: "procedure DoSomething(X: Real = 1.0; I: Integer = 0; S: string = '');"},
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
