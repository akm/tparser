package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestAssemblerStatement(t *testing.T) {
	runBlock := func(t *testing.T, name string, text []rune, expected *ast.Block) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseBlock()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	runBlock(t,
		"without variables",
		[]rune(`
begin
	asm
		DB     FFH                           { One byte }
		DB     0.99                          { Two bytes }
		DB     'A'                           { Ord('A') }
		DB     'Hello world...',0DH,0AH      { String followed by CR/LF }
		DB     12,'string'                   { {{Delphi}} style string }
		DW     0FFFFH                        { One word }
		DW     0,9999                        { Two words }
		DW   'A'                             { Same as DB  'A',0 }
		DW   'BA'                            { Same as DB 'A','B' }
		DW   MyVar                           { Offset of MyVar }
		DW   MyProc                          { Offset of MyProc }
		DD   0FFFFFFFFH                      { One double-word }
		DD   0,999999999         { Two double-words }
		DD   'A'             { Same as DB 'A',0,0,0 }
		DD   'DCBA'              { Same as DB 'A','B','C','D' }
		DD   MyVar               { Pointer to MyVar }
		DD   MyProc              { Pointer to MyProc }
	end;
end;
`),
		&ast.Block{
			CompoundStmt: &ast.CompoundStmt{
				StmtList: ast.StmtList{
					&ast.Statement{
						Body: &ast.AssemblerStatement{},
					},
				},
			},
		},
	)

	runBlock(t,
		"with variables", // but they are not ignored
		[]rune(`
var
	ByteVar: Byte;
	WordVar: Word;
	IntVar: Integer;
asm
	MOV AL,ByteVar
	MOV BX,WordVar
	MOV ECX,IntVar
end;
`),
		&ast.Block{
			DeclSections: ast.DeclSections{
				ast.VarSection{
					&ast.VarDecl{IdentList: asttest.NewIdentList("ByteVar"), Type: asttest.NewOrdIdent("Byte")},
					&ast.VarDecl{IdentList: asttest.NewIdentList("WordVar"), Type: asttest.NewOrdIdent("Word")},
					&ast.VarDecl{IdentList: asttest.NewIdentList("IntVar"), Type: asttest.NewOrdIdent("Integer")},
				},
			},
			CompoundStmt: &ast.CompoundStmt{
				StmtList: ast.StmtList{
					&ast.Statement{
						Body: &ast.AssemblerStatement{},
					},
				},
			},
		},
	)

	runWithFunc := func(t *testing.T, name string, text []rune, expected *ast.FunctionDecl) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseProcedureDeclSection()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	declE := &ast.FormalParm{
		Parameter: &ast.Parameter{
			IdentList: asttest.NewIdentList("e"),
			Type: &ast.ParameterType{
				Type: asttest.NewTypeId("TExample"),
			},
		},
	}
	runWithFunc(t, "procedure with asm block", []rune(`
procedure CallVirtualMethod(e: TExample);
asm
	// Instance pointer needs to be in EAX
	MOV     EAX, e
	// Retrieve VMT table entry
	MOV     EDX, [EAX]
	// Now call the method at offset VMTOFFSET
	CALL    DWORD PTR [EDX + VMTOFFSET TExample.VirtualMethod]
end;
`), &ast.FunctionDecl{
		FunctionHeading: &ast.FunctionHeading{
			Type:             ast.FtProcedure,
			Ident:            asttest.NewIdent("CallVirtualMethod"),
			FormalParameters: ast.FormalParameters{declE},
		},
		Block: &ast.Block{},
	})

}
