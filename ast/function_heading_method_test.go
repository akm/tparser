package ast_test

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestFunctionHeadingMethods(t *testing.T) {
	t.Run("Manipulating", func(t *testing.T) {
		ident := asttest.NewIdent("foo", asttest.NewIdentLocation(1, 10, 9, 13))
		exHeading := &ast.ExportedHeading{
			FunctionHeading: &ast.FunctionHeading{Ident: ident},
		}
		fnHeading := exHeading.FunctionHeading

		assert.Equal(t, "foo", fnHeading.Ident.Name)
		assert.Equal(t, asttest.NewIdentLocation(1, 10, 9, 13), fnHeading.Ident.Location)

		fnHeading.Ident.Location = nil
		assert.Nil(t, fnHeading.Ident.Location)
		assert.Nil(t, exHeading.FunctionHeading.Ident.Location)
	})

	t.Run("Children", func(t *testing.T) {
		ident := asttest.NewIdent("foo", asttest.NewIdentLocation(1, 10, 9, 13))
		fnHeadingOrig := &ast.FunctionHeading{Ident: ident}
		t.Run("*FunctionHeading", func(t *testing.T) {
			fnHeadingBody := *fnHeadingOrig
			fnHeading := &fnHeadingBody
			children := fnHeading.Children()
			assert.Equal(t, 1, len(children))
			assert.Equal(t, ident, children[0])
		})
		t.Run("*ExportedHeading", func(t *testing.T) {
			exHeading := &ast.ExportedHeading{FunctionHeading: fnHeadingOrig}
			children := exHeading.Children()
			assert.Equal(t, 1, len(children))
			assert.Equal(t, fnHeadingOrig, children[0])
			grandChildren := children[0].Children()
			assert.Equal(t, 1, len(grandChildren))
			assert.Equal(t, ident, grandChildren[0])
		})
	})

	t.Run("ClearLocations", func(t *testing.T) {
		ident1 := asttest.NewIdent("Proc1", asttest.NewIdentLocation(1, 11, 10, 16))
		parameters := ast.FormalParameters{
			asttest.NewFormalParm(
				asttest.NewIdent("Param1", asttest.NewIdentLocation(1, 17, 16, 23)),
				asttest.NewOrdIdent(asttest.NewIdent("INTEGER", asttest.NewIdentLocation(1, 25, 24, 32))),
			),
		}

		ident2Body := *ident1
		ident2 := &ident2Body
		asttest.ClearLocations(t, ident2)
		assert.Nil(t, ident2.Location)

		ident3Body := *ident1
		ident3 := &ident3Body
		asttest.ClearLocation(ident3)
		assert.Nil(t, ident3.Location)

		heading := &ast.FunctionHeading{
			Type:             ast.FtProcedure,
			Ident:            ident1,
			FormalParameters: parameters,
		}
		assert.Equal(t, ast.Nodes{ident1, parameters}, heading.Children())

		identNames := []string{}
		astcore.WalkDown(heading, func(node ast.Node) error {
			if ident, ok := node.(*ast.Ident); ok {
				identNames = append(identNames, ident.Name)
			}
			return nil
		})
		assert.Equal(t, []string{"Proc1", "Param1", "INTEGER"}, identNames)

		asttest.ClearLocations(t, heading)

		// ast.WalkDown(heading, func(node ast.Node) error {
		// 	if ident, ok := node.(*ast.Ident); ok {
		// 		// ident.Name = ident.Name + "X"
		// 		ident.Location = nil
		// 		assert.Nil(t, ident.Location, "ident: %s 1", ident.Name)
		// 	}
		// 	return nil
		// })

		assert.Equal(t, &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("Proc1"),
			FormalParameters: ast.FormalParameters{
				asttest.NewFormalParm(
					asttest.NewIdent("Param1"),
					asttest.NewOrdIdent(asttest.NewIdent("INTEGER")),
				),
			},
		}, heading)
	})
}
