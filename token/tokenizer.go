package token

import (
	"github.com/akm/tparser/runes"
)

type TokeninzerFlag uint

const (
	LoadSpace TokeninzerFlag = 1 << iota
	LoadComment
)

type Tokenizer struct {
	*runes.Cursor
	loadSpace   bool
	loadComment bool
}

func NewTokenizer(text *[]rune, flags TokeninzerFlag) *Tokenizer {
	return &Tokenizer{
		Cursor:      runes.NewCuror(text),
		loadSpace:   flags&LoadSpace == LoadSpace,
		loadComment: flags&LoadComment == LoadComment,
	}
}

var processors = []func(*runes.Cursor) *Token{
	ProcessEof,
	ProcessComment,
	ProcessString,
	ProcessNumeral,
	ProcessDoubleSpecialSymbol,
	ProcessSingleSpecialSymbol,
	ProcessWord,
	ProcessSpace,
}

func (t *Tokenizer) GetNext() *Token {
	for _, proc := range processors {
		token := proc(t.Cursor)
		if token != nil {
			if !t.loadSpace && token.Type == Space {
				return t.GetNext()
			} else if !t.loadComment && token.Type == Comment {
				return t.GetNext()
			} else {
				return token
			}
		}
	}
	return nil
}
