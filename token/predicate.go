package token

import (
	"fmt"
	"strings"

	"github.com/akm/tparser/ext"
)

type Predicator interface {
	Name() string
	Predicate(*Token) bool
}

type TokenPredicateImpl struct {
	name      string
	predicate func(*Token) bool
}

func (p *TokenPredicateImpl) Name() string {
	return p.name
}

func (p *TokenPredicateImpl) Predicate(token *Token) bool {
	return p.predicate(token)
}

func OneOf(values ...string) Predicator {
	texts := ext.Strings(values).ToUpper().Set()
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("One of %v", values),
		predicate: func(t *Token) bool { return texts.Include(strings.ToUpper(t.Value())) },
	}
}

func TokenType(typ Type) Predicator {
	return &TokenPredicateImpl{
		name:      typ.String(),
		predicate: func(t *Token) bool { return t.Type == typ },
	}
}

func Symbol(r rune) Predicator {
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("Symbol %q", r),
		predicate: func(t *Token) bool { return t.Type == SpecialSymbol && t.Raw()[0] == r },
	}
}

func Value(v string) Predicator {
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("Value %q", v),
		predicate: func(t *Token) bool { return t.Value() == v },
	}
}

func UpperCase(v string) Predicator {
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("Value %q", v),
		predicate: func(t *Token) bool { return strings.ToUpper(t.Value()) == v },
	}
}
