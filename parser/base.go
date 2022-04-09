package parser

import (
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Parser struct {
	tokenizer *token.Tokenizer

	curr *token.Token
}

func NewParser(text *[]rune) *Parser {
	return &Parser{
		tokenizer: token.NewTokenizer(text, 0),
	}
}

func (p *Parser) RollbackPoint() func() {
	tokenizer := p.tokenizer.Clone()
	curr := p.curr.Clone()
	return func() {
		p.tokenizer = tokenizer
		p.curr = curr
	}
}

func (p *Parser) NextToken() *token.Token {
	p.curr = p.tokenizer.GetNext()
	return p.curr
}

func (p *Parser) CurrentToken() *token.Token {
	return p.curr
}

func (p *Parser) Next(pred token.Predicator) (*token.Token, error) {
	token := p.NextToken()
	if err := p.Validate(token, pred); err != nil {
		return nil, err
	}
	return p.curr, nil
}

func (p *Parser) Current(pred token.Predicator) (*token.Token, error) {
	if err := p.Validate(p.CurrentToken(), pred); err != nil {
		return nil, err
	}
	return p.curr, nil
}

func (p *Parser) Validate(token *token.Token, predicates ...token.Predicator) error {
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

func (p *Parser) Until(terminator token.Predicator, separator token.Predicator, fn func() error) error {
	for {
		if err := fn(); err != nil {
			return err
		}
		token := p.CurrentToken()
		if token == nil {
			return errors.Errorf("something wrong, token is nil")
		}
		if terminator.Predicate(token) {
			break
		}
		if separator != nil {
			if err := p.Validate(token, separator); err != nil {
				return err
			}
			p.NextToken()
		}
	}
	return nil
}
