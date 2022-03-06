package token

import (
	"fmt"

	"github.com/akm/opparser/runes"
)

type Token struct {
	Type  Type
	text  *[]rune
	Start *runes.Position
	End   *runes.Position
}

func NewToken(typ Type, text *[]rune, start, end *runes.Position) *Token {
	return &Token{Type: typ, text: text, Start: start, End: end}
}

func (t *Token) Raw() []rune {
	return (*t.text)[t.Start.Index:t.End.Index]
}

func (t *Token) Text() string {
	return string(t.Raw())
}

func (t *Token) Len() int {
	return t.End.Index - t.Start.Index
}

func (t *Token) TextAbbr(n int) string {
	l := t.Len()
	if l < n {
		return t.Text()
	} else {
		st := t.Start.Index
		return string((*t.text)[st:(st+n)]) + "..."
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %q at %d:%d", t.Type, t.TextAbbr(20), t.Start.Line, t.Start.Col)
}

func (t *Token) Is(pred TokenPredicate) bool {
	return pred.Predicate(t)
}
