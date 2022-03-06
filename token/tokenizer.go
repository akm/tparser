package token

import (
	"github.com/akm/opparser/runes"
)

type Tokenizer struct {
	*runes.Cursor
}

func NewTokenizer(text *[]rune) *Tokenizer {
	return &Tokenizer{
		Cursor: runes.NewCuror(text),
	}
}

var processors = []func(*runes.Cursor) *Token{
	ProcessEof,
	ProcessComment,
	ProcessString,
	ProcessNumeral,
	ProcessSingleSpecialSymbol,
	ProcessWord,
	ProcessSpace,
}

func (t *Tokenizer) Next() *Token {
	for _, proc := range processors {
		if token := proc(t.Cursor); token != nil {
			return token
		}
	}
	return nil
}
