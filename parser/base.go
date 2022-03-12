package parser

import (
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
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

func (p *Parser) NextToken() *token.Token {
	p.prior, p.curr = p.curr, p.tokenizer.GetNext()
	return p.curr
}

func (p *Parser) CurrentToken() *token.Token {
	return p.curr
}

func (p *Parser) Next(pred token.Predicate) (*token.Token, error) {
	token := p.NextToken()
	if err := p.Validate(token); err != nil {
		return nil, err
	}
	return p.curr, nil
}

func (p *Parser) Current(pred token.Predicate) (*token.Token, error) {
	if err := p.Validate(p.CurrentToken()); err != nil {
		return nil, err
	}
	return p.curr, nil
}

func (p *Parser) Validate(token *token.Token, predicates ...token.Predicate) error {
	if token == nil {
		return errors.Errorf("something wrong, token is nil")
	}
	for _, pred := range predicates {
		if !pred.Predicate(token) {
			return errors.Errorf("expects %s but was %s", pred.Name(), token.String())
		}
	}
	return nil
}

func (p *Parser) Until(terminator token.Predicate, separator token.Predicate, fn func() error) error {
	for {
		if err := fn(); err != nil {
			return err
		}
		token := p.NextToken()
		if token == nil {
			return errors.Errorf("something wrong, token is nil")
		}
		if terminator.Predicate(token) {
			break
		}
		if separator != nil {
			if p.Validate(token, separator) != nil {
				return errors.Errorf("expects %s but was %s", separator.Name(), token.String())
			}
		}
	}
	return nil
}
