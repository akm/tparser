package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestStrucFileType(t *testing.T) {
	declPhoneEntry := &ast.TypeDecl{
		Ident: asttest.NewIdent("PhoneEntry"),
		Type: &ast.RecType{
			FieldList: &ast.FieldList{
				FieldDecls: ast.FieldDecls{
					{
						IdentList: asttest.NewIdentList("FirstName", "LastName"),
						Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("20"))},
					},
					{
						IdentList: asttest.NewIdentList("PhoneNumber"),
						Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("15"))},
					},
					{
						IdentList: asttest.NewIdentList("Listed"),
						Type:      &ast.OrdIdent{Ident: asttest.NewIdent("Boolean")},
					},
				},
			},
		},
	}

	NewTypeSectionTestRunner(t,
		"with TShapeList",
		[]rune(`
type
	PhoneEntry = record
		FirstName, LastName: string[20];
		PhoneNumber: string[15];
		Listed: Boolean;
	end;
	PhoneList = file of PhoneEntry;
`),
		ast.TypeSection{
			declPhoneEntry,
			{
				Ident: asttest.NewIdent("PhoneList"),
				Type: &ast.FileType{
					TypeId: asttest.NewTypeId("PhoneEntry", declPhoneEntry.ToDeclarations()[0]),
				},
			},
		},
	).Run()
}
