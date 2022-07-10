package astcore

import (
	"fmt"
	"strings"

	"github.com/akm/tparser/runes"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Position = runes.Position
type Location struct {
	Path  string
	Start *Position
	End   *Position
}

func NewLocation(start, end *runes.Position) *Location {
	return &Location{Start: start, End: end}
}

func (m *Location) String() string {
	return fmt.Sprintf("%s(%s - %s)", m.Path, m.Start, m.End)
}

// - Ident
//   ```
//   <identifier>
//   ```
type Ident struct {
	Name     string
	Location *Location
}

func NewIdent(v *token.Token) *Ident {
	return &Ident{
		Name:     v.RawString(),
		Location: NewLocation(v.Start, v.End),
	}
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

func (m *Ident) Children() Nodes {
	return Nodes{}
}

func (m *Ident) String() string {
	if m == nil {
		return ""
	} else {
		return m.Name
	}
}

// - IdentList
//   ```
//   Ident ','...
//   ```
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

func (s IdentList) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

func (s IdentList) Find(name string) *Ident {
	kw := strings.ToLower(name)
	for _, i := range s {
		if strings.ToLower(i.Name) == kw {
			return i
		}
	}
	return nil
}
