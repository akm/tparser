package parser

import (
	"github.com/akm/opparser/token"
)

type Parser struct {
	tokenizer *token.Tokenizer

	curr  *token.Token
	prior *token.Token
}

func NewParser(text *[]rune) *Parser {
	return &Parser{
		tokenizer: token.NewTokenizer(text, 0),
	}
}

func (p *Parser) next() *token.Token {
	p.prior, p.curr = p.curr, p.tokenizer.Next()
	return p.curr
}

func (p *Parser) get(pred token.Predicate) (*token.Token, error) {
	token, err := p.tokenizer.Get(pred)
	if err != nil {
		return nil, err
	}
	p.prior, p.curr = p.curr, token
	return p.curr, nil
}
