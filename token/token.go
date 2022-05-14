package token

import (
	"fmt"
	"strings"

	"github.com/akm/tparser/runes"
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

func (t *Token) Clone() *Token {
	return &Token{
		Type:  t.Type,
		text:  t.text,
		Start: t.Start.Clone(),
		End:   t.End.Clone(),
	}
}

func (t *Token) Raw() []rune {
	return (*t.text)[t.Start.Index:t.End.Index]
}

func (t *Token) RawString() string {
	return string(t.Raw())
}

func (t *Token) Value() string {
	res := t.RawString()
	switch t.Type {
	case ReservedWord:
		return strings.ToUpper(res)
	default:
		return res
	}
}

func (t *Token) Len() int {
	return t.End.Index - t.Start.Index
}

func (t *Token) ValueAbbr(n int) string {
	l := t.Len()
	if l < n {
		return t.Value()
	} else {
		st := t.Start.Index
		return string((*t.text)[st:(st+n)]) + "..."
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %q at %d:%d", t.Type, t.ValueAbbr(20), t.Start.Line, t.Start.Col)
}

func (t *Token) Is(pred Predicator) bool {
	return pred.Predicate(t)
}
