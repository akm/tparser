package parser

import (
	"fmt"

	"github.com/akm/tparser/log"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Parser struct {
	tokenizer        *token.Tokenizer
	curr             *token.Token
	context          Context
	postSectionFuncs []func()
}

func NewParser(ctx Context) *Parser {
	if ctx == nil {
		panic(errors.Errorf("context is required for NewParser"))
	}
	return &Parser{context: ctx}
}

func (p *Parser) SetText(text *[]rune) {
	p.tokenizer = token.NewTokenizer(text, 0)
}

func (p *Parser) RollbackPoint() func() {
	tokenizer := p.tokenizer.Clone()
	curr := p.curr.Clone()
	ctx := p.context.Clone()
	return func() {
		p.tokenizer = tokenizer
		p.curr = curr
		p.context = ctx
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

func (p *Parser) Validate(t *token.Token, predicates ...token.Predicator) error {
	if t == nil {
		return errors.Errorf("something wrong, token is nil")
	}
	for _, pred := range predicates {
		if !pred.Predicate(t) {
			return p.TokenErrorf("expects "+pred.Name()+" but was %s", t)
		}
	}
	return nil
}

type QuitUntilType struct{}

func (*QuitUntilType) Error() string {
	return "quit until"
}

var QuitUntil = &QuitUntilType{}

func (p *Parser) Until(terminator token.Predicator, separator token.Predicator, fn func() error) error {
	for {
		if err := fn(); err != nil {
			if err == QuitUntil {
				return nil
			}
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

func (p *Parser) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (p *Parser) TraceMethod(name string) func() {
	p.Logf("%s START %s\n", name, p.CurrentToken().String())
	return func() { p.Logf("%s END   %s\n", name, p.CurrentToken().String()) }
}

func (p *Parser) TokenErrorf(format string, t *token.Token, args ...interface{}) error {
	fmtArgs := append([]interface{}{t.RawString()}, args...)
	fmtArgs = append(fmtArgs, p.PlaceString(t))
	return errors.Errorf(format+" at %s", fmtArgs...)
}

func (p *Parser) PlaceString(t *token.Token) string {
	return fmt.Sprintf("%s:%d:%d", p.context.GetPath(), t.Start.Line, t.Start.Col)
}

func (p *Parser) SetupPostSectionFuncs() func() {
	var backup []func()
	p.postSectionFuncs, backup = []func(){}, p.postSectionFuncs
	return func() { p.postSectionFuncs = backup }
}

func (p *Parser) AddPostSection(fn func()) {
	p.postSectionFuncs = append(p.postSectionFuncs, fn)
}

func (p *Parser) RunPostSectionFuncs() {
	for _, fn := range p.postSectionFuncs {
		fn()
	}
}
