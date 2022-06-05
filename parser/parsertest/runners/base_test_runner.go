package runners

import (
	"testing"

	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

type BaseTestRunner struct {
	T              *testing.T
	Name           string
	Text           *[]rune
	ClearLocations bool
	ParserArgFuncs ParserArgFuncs
}

func NewBaseTestRunner(
	t *testing.T,
	name string,
	text *[]rune,
	clearLocations bool,
	parserArgFuncs ParserArgFuncs,
	runnerFuncs BaseTestRunnerFuncs,
) *BaseTestRunner {
	r := &BaseTestRunner{
		T:              t,
		Name:           name,
		Text:           text,
		ClearLocations: clearLocations,
		ParserArgFuncs: parserArgFuncs,
	}
	for _, f := range runnerFuncs {
		f(r)
	}
	return r
}

func (tt *BaseTestRunner) NewParser() *parser.Parser {
	args := make([]interface{}, len(tt.ParserArgFuncs))
	for i, f := range tt.ParserArgFuncs {
		args[i] = f()
	}
	r := NewTestParser(tt.Text, args...)
	r.NextToken()
	return r
}

func (tt *BaseTestRunner) Run(
	parseFunc func(*parser.Parser) (astcore.Node, error),
	assertFunc func(t *testing.T, actual astcore.Node),
) {
	tt.T.Run(tt.Name, func(t *testing.T) {
		res, err := parseFunc(tt.NewParser())
		if assert.NoError(t, err) {
			if tt.ClearLocations {
				asttest.ClearLocations(t, res)
			}
			assertFunc(t, res)
		}
	})
}

type ParserArgFunc = func() interface{}
type ParserArgFuncs []ParserArgFunc

func FilterParserArgFuncs(args []interface{}) (ParserArgFuncs, []interface{}) {
	r := ParserArgFuncs{}
	rest := []interface{}{}
	for _, arg := range args {
		if v, ok := arg.(ParserArgFunc); ok {
			r = append(r, v)
		} else {
			rest = append(rest, arg)
		}
	}
	return r, rest
}

func (s ParserArgFuncs) Interfaces() []interface{} {
	r := make([]interface{}, len(s))
	for i, f := range s {
		r[i] = f
	}
	return r
}

type BaseTestRunnerFunc = func(*BaseTestRunner)
type BaseTestRunnerFuncs []BaseTestRunnerFunc

func FilterBaseTestRunnerFuncs(args []interface{}) (BaseTestRunnerFuncs, []interface{}) {
	r := BaseTestRunnerFuncs{}
	rest := []interface{}{}
	for _, arg := range args {
		if v, ok := arg.(BaseTestRunnerFunc); ok {
			r = append(r, v)
		} else {
			rest = append(rest, arg)
		}
	}
	return r, rest
}
