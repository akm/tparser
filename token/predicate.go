package token

import (
	"fmt"
	"strings"
)

type Predicator interface {
	Name() string
	Predicate(*Token) bool
}

type Predicators []Predicator

func (s Predicators) Names() []string {
	names := make([]string, len(s))
	for i, p := range s {
		names[i] = p.Name()
	}
	return names
}

func (s Predicators) Some(t *Token) bool {
	for _, p := range s {
		if p.Predicate(t) {
			return true
		}
	}
	return false
}

func (s Predicators) Every(t *Token) bool {
	for _, p := range s {
		if !p.Predicate(t) {
			return false
		}
	}
	return true
}

type PredicatorImpl struct {
	name      string
	predicate func(*Token) bool
}

func (p *PredicatorImpl) Name() string {
	return p.name
}

func (p *PredicatorImpl) Predicate(token *Token) bool {
	return p.predicate(token)
}

func Some(predicators ...Predicator) Predicator {
	s := Predicators(predicators)
	return &PredicatorImpl{
		name:      fmt.Sprintf("Some of %v", s.Names()),
		predicate: s.Some,
	}
}

func Every(predicators ...Predicator) Predicator {
	s := Predicators(predicators)
	return &PredicatorImpl{
		name:      fmt.Sprintf("Every of %v", s.Names()),
		predicate: s.Every,
	}
}

func TokenType(typ Type) Predicator {
	return &PredicatorImpl{
		name:      typ.String(),
		predicate: func(t *Token) bool { return t.Type == typ },
	}
}

func Symbol(r rune) Predicator {
	return &PredicatorImpl{
		name:      fmt.Sprintf("Symbol %q", r),
		predicate: func(t *Token) bool { return t.Type == SpecialSymbol && t.Raw()[0] == r },
	}
}

func Value(v string) Predicator {
	return &PredicatorImpl{
		name:      fmt.Sprintf("Value %q", v),
		predicate: func(t *Token) bool { return t.Value() == v },
	}
}

func UpperCase(v string) Predicator {
	return &PredicatorImpl{
		name:      fmt.Sprintf("Value %q", v),
		predicate: func(t *Token) bool { return strings.ToUpper(t.Value()) == v },
	}
}
