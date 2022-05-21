package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Program struct {
	*ast.Program
	Units ast.Units
}

func ParseProgram(path string) (*Program, error) {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	decoder := japanese.ShiftJIS.NewDecoder()
	str, err := ioutil.ReadAll(transform.NewReader(fp, decoder))
	if err != nil {
		return nil, err
	}

	runes := []rune(string(str))

	// absPath, err := filepath.Abs(path)
	// if err != nil {
	// 	return nil, err
	// }

	ctx := NewProgramContext(path)
	p := NewParser(&runes, ctx)
	p.NextToken()
	res, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	return &Program{
		Program: res,
		Units:   ctx.Units,
	}, nil
}

type Parser struct {
	tokenizer *token.Tokenizer
	curr      *token.Token
	context   Context
	logger    *log.Logger
}

func NewParser(text *[]rune, ctx Context, args ...interface{}) *Parser {
	var logger *log.Logger
	for _, arg := range args {
		switch v := arg.(type) {
		case *log.Logger:
			logger = v
		default:
			panic(errors.Errorf("unexpected type %T (%v)", arg, arg))
		}
	}
	if ctx == nil {
		panic(errors.Errorf("context is required for NewParser"))
	}
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)
	}
	return &Parser{
		tokenizer: token.NewTokenizer(text, 0),
		context:   ctx,
		logger:    logger,
	}
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

func (p *Parser) Logf(format string, args ...interface{}) {
	p.logger.Printf(format, args...)
}

func (p *Parser) TokenErrorf(format string, t *token.Token, args ...interface{}) error {
	fmtArgs := append([]interface{}{t.RawString()}, args...)
	place := fmt.Sprintf("%s:%d,%d", p.context.GetPath(), t.Start.Line, t.Start.Col)
	fmtArgs = append(fmtArgs, place)
	return errors.Errorf(format+" at %s", fmtArgs...)
}

func (p *Parser) StackContext() func() {
	var backup Context
	p.context, backup = NewStackableContext(p.context), p.context
	return func() {
		p.context = backup
	}
}
