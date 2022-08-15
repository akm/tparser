package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/log/testlog"
)

func TestClassRefType(t *testing.T) {
	defer testlog.Setup(t)()

	RunTypeSection(t,
		"class ref without class",
		[]rune(`type TClass = class of TObject;`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: asttest.NewIdent("TClass"),
				Type:  ast.NewCustomClassRefType(asttest.NewTypeId("TObject")),
			},
		},
	)

	RunTypeSection(t,
		"class ref with class",
		[]rune(`
type
	TFigure = class
	end;
    TFigureClass = class of TFigure;
`),
		func() ast.TypeSection {
			classDecl := &ast.TypeDecl{
				Ident: asttest.NewIdent("TFigure"),
				Type:  &ast.CustomClassType{},
			}
			return ast.TypeSection{
				classDecl,
				&ast.TypeDecl{
					Ident: asttest.NewIdent("TFigureClass"),
					Type:  ast.NewCustomClassRefType(asttest.NewTypeId("TFigure", classDecl.ToDeclarations()[0])),
				},
			}
		}(),
	)
}
