package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func NewIdent(arg interface{}, locations ...*ast.Location) *ast.Ident {
	var r *ast.Ident
	switch v := arg.(type) {
	case string:
		r = &ast.Ident{Name: v}
	case *string:
		r = &ast.Ident{Name: *v}
	default:
		r = ast.NewIdentFrom(arg)
	}
	if len(locations) > 0 {
		r.Location = locations[0]
	}
	return r
}

func NewIdentList(args ...interface{}) ast.IdentList {
	switch len(args) {
	case 0:
		panic(errors.Errorf("no arguments are given for NewIdentList"))
	case 1:
		arg := args[0]
		switch v := arg.(type) {
		case string:
			return ast.NewIdentList(NewIdent(v))
		case *string:
			return ast.NewIdentList(NewIdent(v))
		case []string:
			r := make(ast.IdentList, len(v))
			for idx, i := range v {
				r[idx] = NewIdent(i)
			}
			return r
		default:
			return ast.NewIdentList(arg)
		}
	default:
		r := make(ast.IdentList, len(args))
		for i, arg := range args {
			r[i] = NewIdent(arg)
		}
		return r
	}
}

// NewPosition(line, col, index)
// NewPosition(1, col)
// NewPosition(col) // line is 1
func NewPosition(args ...int) *ast.Position {
	switch len(args) {
	case 1:
		col := args[0]
		return NewPosition(1, col)
	case 2:
		line, col := args[0], args[1]
		if line == 1 {
			return NewPosition(line, col, col+1)
		} else {
			panic(errors.Errorf("index required unless line is not 1, line: %d, col: %d", line, col))
		}
	case 3:
		line, col, index := args[0], args[1], args[2]
		return &ast.Position{Line: line, Col: col, Index: index}
	default:
		panic(errors.Errorf("unexpected number of arguments (%d) are given for NewPosition", len(args)))
	}
}

func castPosition(arg interface{}) *ast.Position {
	switch v := arg.(type) {
	case ast.Position:
		return &v
	case *ast.Position:
		return v
	default:
		panic(errors.Errorf("unexpected type %T (%v) is given for Position", arg, arg))
	}
}
func castInt(arg interface{}) int {
	switch v := arg.(type) {
	case int:
		return v
	case *int:
		return *v
	default:
		panic(errors.Errorf("unexpected type %T (%v) is given for int", arg, arg))
	}
}

func castInts(args []interface{}) []int {
	r := make([]int, len(args))
	for i, v := range args {
		r[i] = castInt(v)
	}
	return r
}

// args are
// - *token.Token
// - *ast.Position, *ast.Position
// - startLine, startCol, startIndex, endCol
// - startLine, startCol, startIndex, endCol, endIndex
// - startLine, startCol, startIndex, endLine, endCol, endIndex
func NewIdentLocation(args ...interface{}) *ast.Location {
	switch len(args) {
	case 1:
		switch v := args[0].(type) {
		case ast.Location:
			return &v
		case *ast.Location:
			return v
		case token.Token:
			return ast.NewLocation(v.Start, v.End)
		case *token.Token:
			return ast.NewLocation(v.Start, v.End)
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewIdentLocation", args[0], args[0]))
		}
	case 2:
		return &ast.Location{Start: castPosition(args[0]), End: castPosition(args[1])}
	case 4:
		vals := castInts(args)
		startLine, startCol, startIndex, endCol := vals[0], vals[1], vals[2], vals[3]
		return NewIdentLocation(startLine, startCol, startIndex, startLine, endCol, startIndex+(endCol-startCol))
	case 5:
		vals := castInts(args)
		startLine, startCol, startIndex, endCol, endIndex := vals[0], vals[1], vals[2], vals[3], vals[4]
		if (endCol - startCol) != (endIndex - startIndex) {
			panic(errors.Errorf("length conflicted: startIndex: %d, startCol: %d, endIndex: %d, endCol: %d", startIndex, startCol, endIndex, endCol))
		}
		return NewIdentLocation(startLine, startCol, startIndex, startLine, endCol, endIndex)
	case 6:
		vals := castInts(args)
		return &ast.Location{Start: NewPosition(vals[0:3]...), End: NewPosition(vals[3:6]...)}
	default:
		panic(errors.Errorf("unexpected number of arguments (%d) are given for NewIdentLocation", len(args)))
	}
}

func ClearLocation(ident *ast.Ident) {
	ident.Location = nil
}

func ClearLocations(t *testing.T, node ast.Node) {
	err := astcore.WalkDown(node, func(n ast.Node) error {
		switch v := n.(type) {
		case *ast.Ident:
			ClearLocation(v)
		}
		return nil
	})
	assert.NoError(t, err)
}
