package token

import (
	"github.com/akm/opparser/runes"
	"github.com/pkg/errors"
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
	ProcessSingleSpecialSymbol,
	ProcessWord,
	ProcessSpace,
}

func (t *Tokenizer) Next() *Token {
	for _, proc := range processors {
		token := proc(t.Cursor)
		if token != nil {
			if !t.loadSpace && token.Type == Space {
				return t.Next()
			} else if !t.loadComment && token.Type == Comment {
				return t.Next()
			} else {
				return token
			}
		}
	}
	return nil
}

func (t *Tokenizer) Get(pred TokenPredicate) (*Token, error) {
	token := t.Next()
	if token == nil {
		return nil, errors.Errorf("something wrong, token is nil")
	}
	if !pred.Predicate(token) {
		return nil, errors.Errorf("expects %s but was %s", pred.Name(), token.String())
	}
	return token, nil
}
