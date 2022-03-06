package token

import (
	"fmt"
	"strings"

	"github.com/akm/opparser/ext"
)

type TokenPredicate interface {
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

func OneOf(values ...string) TokenPredicate {
	texts := ext.Strings(values).ToUpper().Set()
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("One of %v", values),
		predicate: func(t *Token) bool { return texts.Include(strings.ToUpper(t.Text())) },
	}
}

func TokenType(typ Type) TokenPredicate {
	return &TokenPredicateImpl{
		name:      typ.String(),
		predicate: func(t *Token) bool { return t.Type == typ },
	}
}

func Symbol(r rune) TokenPredicate {
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("Symbol %q", r),
		predicate: func(t *Token) bool { return t.Type == SpecialSymbol && t.Raw()[0] == r },
	}
}
