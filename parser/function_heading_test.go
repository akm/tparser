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
	type pattern struct {
		text     string
		expected *ast.ExportedHeading
	}

	patterns := []pattern{}

	run := func(text string, expected *ast.ExportedHeading) {
		t.Run(text, func(t *testing.T) {
			patterns = append(patterns, pattern{text, expected})
			runes := []rune(text)
			parser := NewParser(&runes, NewContext())
			parser.NextToken()
			res, err := parser.ParseExportedHeading()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		`PROCEDURE Proc0;`,
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("Proc0"),
			},
		},
	)

	run(
		`PROCEDURE Proc1(Param1: INTEGER);`,
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("Proc1"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("Param1", ast.NewOrdIdent("INTEGER")),
				},
			},
		},
	)

	run(
		`procedure NumString(N: Integer; var S: string);`,
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("NumString"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("N", ast.NewOrdIdent("Integer")),
					ast.NewFormalParm("S", ast.NewStringType("STRING"), &ast.FpoVar),
				},
			},
		},
	)

	run(
		"function WF: Integer;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:       ast.FtFunction,
				Ident:      ast.NewIdent("WF"),
				ReturnType: ast.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"function Max(A: array of Real; N: Integer): Real;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("Max"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("A", ast.NewArrayParameterType(ast.NewRealType("Real"))),
					ast.NewFormalParm("N", ast.NewOrdIdent("Integer")),
				},
				ReturnType: ast.NewRealType("Real"),
			},
		},
	)

	run(
		"function Power(X: Real; Y: Integer): Real;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("Power"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("X", ast.NewRealType("Real")),
					ast.NewFormalParm("Y", ast.NewOrdIdent("Integer")),
				},
				ReturnType: ast.NewRealType("Real"),
			},
		},
	)

	run(
		"function MyFunction(X, Y: Real): Real; cdecl;",
		&ast.ExportedHeading{
			Directives: []ast.Directive{ast.DrCdecl},
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("MyFunction"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"X", "Y"}, ast.NewRealType("Real")),
				},
				ReturnType: ast.NewRealType("Real"),
			},
		},
	)

	run(
		"function Calculate(X, Y: Integer): Real; forward;",
		&ast.ExportedHeading{
			Directives: []ast.Directive{ast.DrForward},
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("Calculate"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"X", "Y"}, ast.NewOrdIdent("Integer")),
				},
				ReturnType: ast.NewRealType("Real"),
			},
		},
	)

	run(
		"function printf(Format: PChar): Integer; cdecl; varargs;",
		&ast.ExportedHeading{
			Directives: []ast.Directive{ast.DrCdecl, ast.DrVarArgs},
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("printf"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("Format", "PChar"),
				},
				ReturnType: ast.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"procedure MoveWord(var Source, Dest; Count: Integer); external;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("MoveWord"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
					ast.NewFormalParm("Count", ast.NewOrdIdent("Integer")),
				},
			},
			Directives: []ast.Directive{ast.DrExternal},
		},
	)

	run(
		"function SomeFunction(S: string): string; external 'strlib.dll';",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("SomeFunction"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("S", ast.NewStringType("STRING")),
				},
				ReturnType: ast.NewStringType("STRING"),
			},
			Directives:      []ast.Directive{ast.DrExternal},
			ExternalOptions: &ast.ExternalOptions{LibraryName: "'strlib.dll'"},
		},
	)

	run(
		"function MessageBox(HWnd: Integer; Text, Caption: PChar; Flags: Integer): Integer; stdcall; external 'user32.dll' name 'MessageBoxA';",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("MessageBox"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm("HWnd", ast.NewOrdIdent("Integer")),
					ast.NewFormalParm([]string{"Text", "Caption"}, "PChar"),
					ast.NewFormalParm("Flags", ast.NewOrdIdent("Integer")),
				},
				ReturnType: ast.NewOrdIdent("Integer"),
			},
			Directives:      []ast.Directive{ast.DrStdcall, ast.DrExternal},
			ExternalOptions: &ast.ExternalOptions{LibraryName: "'user32.dll'", Name: ext.StringPtr("'MessageBoxA'")},
		},
	)

	run(
		"function Divide(X, Y: Real): Real; overload;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("Divide"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"X", "Y"}, ast.NewRealType("Real")),
				},
				ReturnType: ast.NewRealType("Real"),
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
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("DoubleByValue"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"X"}, ast.NewOrdIdent("Integer")),
				},
				ReturnType: ast.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"function DoubleByRef(var X: Integer): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("DoubleByRef"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"X"}, ast.NewOrdIdent("Integer"), &ast.FpoVar),
				},
				ReturnType: ast.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"function CompareStr(const S1, S2: string): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("CompareStr"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"S1", "S2"}, ast.NewStringType("STRING"), &ast.FpoConst),
				},
				ReturnType: ast.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"procedure GetInfo(out Info: SomeRecordType);",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("GetInfo"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"Info"}, "SomeRecordType", &ast.FpoOut),
				},
			},
		},
	)

	run(
		"function Equal(var Source, Dest; Size: Integer): Boolean;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("Equal"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"Source", "Dest"}, nil, &ast.FpoVar),
					ast.NewFormalParm([]string{"Size"}, ast.NewOrdIdent("Integer")),
				},
				ReturnType: ast.NewOrdIdent("Boolean"),
			},
		},
	)

	// {text: "procedure Check(S: string[20]);"} // syntax error
	run(
		"procedure Check(S: OpenString);",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("Check"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"S"}, "OpenString"),
				},
			},
		},
	)

	// {text: "procedure Sort(A: array[1..10] of Integer);"} // syntax error
	run(
		"function Find(A: array of Char): Integer;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("Find"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"A"}, ast.NewArrayParameterType(ast.NewOrdIdent("Char"))),
				},
				ReturnType: ast.NewOrdIdent("Integer"),
			},
		},
	)

	run(
		"procedure Clear(var A: array of Real);",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("Clear"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"A"}, ast.NewArrayParameterType(ast.NewRealType("Real")), &ast.FpoVar),
				},
			},
		},
	)

	run(
		"function MakeStr(const Args: array of const): string;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("MakeStr"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"Args"}, ast.NewArrayParameterType(nil), &ast.FpoConst),
				},
				ReturnType: ast.NewStringType("STRING"),
			},
		},
	)

	run(
		"procedure FillArray(A: array of Integer; Value: Integer = 0);",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("FillArray"),
				FormalParameters: ast.FormalParameters{
					ast.NewFormalParm([]string{"A"}, ast.NewArrayParameterType(ast.NewOrdIdent("Integer"))),
					ast.NewFormalParm(
						ast.NewParameter([]string{"Value"}, ast.NewOrdIdent("Integer"), ast.NewNumber("0")),
					),
				},
			},
		},
	)

	run(
		"function MyFunction(X: Real = 3.5; Y: Real = 3.5): Real;",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtFunction,
				Ident: ast.NewIdent("MyFunction"),
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
	)

	run(
		"procedure DoSomething(X: Real = 1.0; I: Integer = 0; S: string = '');",
		&ast.ExportedHeading{
			FunctionHeading: ast.FunctionHeading{
				Type:  ast.FtProcedure,
				Ident: ast.NewIdent("DoSomething"),
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
	)
	// {text: "function MyFunction(X, Y: Real = 3.5): Real;"}, // syntax error
	// {text: "procedure MyProcedure(I: Integer = 1; S: string);"} // syntax error

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
					Ident: ast.NewIdent("Unit1"),
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
