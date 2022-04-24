package ast

import (
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Ident struct {
	Name string
}

func NewIdent(arg interface{}) *Ident {
	switch v := arg.(type) {
	case Ident:
		return &v
	case *Ident:
		return v
	case string:
		return &Ident{Name: v}
	case *string:
		return &Ident{Name: *v}
	case token.Token:
		return NewIdent(v.Value())
	case *token.Token:
		return NewIdent(v.Value())
	default:
		panic(errors.Errorf("unexpected type %T (%v) is given for NewIdent", arg, arg))
	}
}

type IdentList []*Ident

func NewIdentList(args ...interface{}) IdentList {
	switch len(args) {
	case 0:
		panic(errors.Errorf("no arguments are given for NewIdentList"))
	case 1:
		arg := args[0]
		switch v := arg.(type) {
		case IdentList:
			return v
		case string:
			return IdentList{NewIdent(v)}
		case *string:
			return IdentList{NewIdent(v)}
		case []string:
			r := make(IdentList, len(v))
			for idx, i := range v {
				r[idx] = NewIdent(i)
			}
			return r
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewIdentList", arg, arg))
		}
	default:
		r := make(IdentList, len(args))
		for i, arg := range args {
			r[i] = NewIdent(arg)
		}
		return r
	}
}

func (s IdentList) Names() []string {
	r := make([]string, len(s))
	for idx, i := range s {
		r[idx] = i.Name
	}
	return r
}
