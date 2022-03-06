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
	l := len(*t.text)
	if t.End.Index < l {
		return (*t.text)[t.Start.Index:t.End.Index]
	} else {
		return (*t.text)[t.Start.Index:]
	}
}

func (t *Token) Text() string {
	return string(t.Raw())
}

func (t *Token) String() string {
	return fmt.Sprintf("%s at %d:%d", t.Type, t.Start.Line, t.Start.Col)
}
