package ast

import (
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Ident struct {
	Name string
}

func NewIdent(v *token.Token) *Ident {
	return &Ident{Name: v.RawString()}
}
func NewIdentFrom(arg interface{}) *Ident {
	switch v := arg.(type) {
	case Ident:
		return &v
	case *Ident:
		return v
	case token.Token:
		return NewIdent(&v)
	case *token.Token:
		return NewIdent(v)
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
		case []*Ident:
			return IdentList(v)
		case []interface{}:
			r := make(IdentList, len(v))
			for idx, i := range v {
				r[idx] = NewIdentFrom(i)
			}
			return r
		case Ident:
			return IdentList{&v}
		case *Ident:
			return IdentList{v}
		case []string:
			r := make(IdentList, len(v))
			for idx, i := range v {
				r[idx] = NewIdentFrom(i)
			}
			return r
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewIdentList", arg, arg))
		}
	default:
		r := make(IdentList, len(args))
		for i, arg := range args {
			r[i] = NewIdentFrom(arg)
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
