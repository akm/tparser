package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/stretchr/testify/assert"
)

func NewUnitRef(name interface{}, paths ...string) *ast.UsesClauseItem {
	switch v := name.(type) {
	case string:
		return ast.NewUnitRef(NewIdent(v), paths...)
	default:
		return ast.NewUnitRef(name, paths...)
	}
}

func ClearUsesItemUnit(item *ast.UsesClauseItem) {
	item.Unit = nil
}

func ClearUsesItemsUnit(t *testing.T, node ast.Node) {
	err := astcore.WalkDown(node, func(n ast.Node) error {
		switch v := n.(type) {
		case *ast.UsesClauseItem:
			ClearUsesItemUnit(v)
		}
		return nil
	})
	assert.NoError(t, err)
}
