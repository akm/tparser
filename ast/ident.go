package ast

import (
	"github.com/akm/tparser/ext"
	"github.com/pkg/errors"
)

type Ident string

func NewIdent(arg interface{}) Ident {
	switch v := arg.(type) {
	case Ident:
		return v
	case *Ident:
		return *v
	case string:
		return Ident(v)
	case *string:
		return Ident(*v)
	default:
		panic(errors.Errorf("unexpected type %T (%v) is given for NewIdent", arg, arg))
	}
}

type IdentList = ext.Strings

func NewIdentList(arg interface{}) IdentList {
	switch v := arg.(type) {
	case IdentList:
		return v
	case string:
		return IdentList{v}
	case *string:
		return IdentList{*v}
	case []string:
		return IdentList(v)
	default:
		panic(errors.Errorf("unexpected type %T (%v) is given for NewIdentList", arg, arg))
	}
}
