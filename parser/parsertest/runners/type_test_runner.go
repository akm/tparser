package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type TypeTestRunner struct {
	t              *testing.T
	Name           string
	Text           *[]rune
	Expected       ast.Type
	ClearLocations bool
	ParserArgFuncs []func() interface{}
	RunnerFuncs    []TypeTestRunnerFunc
}

type TypeTestRunnerFunc = func(*TypeTestRunner)

func NewTypeTestRunner(t *testing.T, name string, text []rune, expected ast.Type, args ...interface{}) *TypeTestRunner {
	parserArgFuncs := []func() interface{}{}
	runnerFuncs := []TypeTestRunnerFunc{}
	for _, arg := range args {
		switch v := arg.(type) {
		case func() interface{}:
			parserArgFuncs = append(parserArgFuncs, v)
		case TypeTestRunnerFunc:
			runnerFuncs = append(runnerFuncs, v)
		default:
			panic(errors.Errorf("unexpected argument type %T %v", v, v))
		}
	}
	return &TypeTestRunner{
		t:              t,
		Name:           name,
		Text:           &text,
		Expected:       expected,
		ParserArgFuncs: parserArgFuncs,
		RunnerFuncs:    runnerFuncs,
		ClearLocations: true,
	}
}

func (tt *TypeTestRunner) newParser(text *[]rune) *parser.Parser {
	for _, fn := range tt.RunnerFuncs {
		fn(tt)
	}
	args := make([]interface{}, len(tt.ParserArgFuncs))
	for i, f := range tt.ParserArgFuncs {
		args[i] = f()
	}
	r := NewTestParser(text, args...)
	r.NextToken()
	return r
}

func (tt *TypeTestRunner) Run() *TypeTestRunner {
	tt.t.Run(tt.Name, func(t *testing.T) {
		p := tt.newParser(tt.Text)
		res, err := p.ParseType()
		if assert.NoError(t, err) {
			if tt.ClearLocations {
				asttest.ClearLocations(t, res)
			}
			assert.Equal(t, tt.Expected, res)
		}
	})
	return tt
}

func (tt *TypeTestRunner) RunTypeSection(declName string) *TypeTestRunner {
	tt.t.Run(tt.Name+" in type section", func(t *testing.T) {
		expectedSection := ast.TypeSection{{Ident: asttest.NewIdent(declName), Type: tt.Expected}}
		sectionStr := "type " + declName + " = " + string(*tt.Text) + ";"
		sectionRunes := []rune(sectionStr)
		p := tt.newParser(&sectionRunes)
		res, err := p.ParseTypeSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expectedSection, res)
		}
	})
	return tt
}

func (tt *TypeTestRunner) RunVarSection(declName string) *TypeTestRunner {
	tt.t.Run(tt.Name+" in var section", func(t *testing.T) {
		expectedSection := ast.VarSection{{IdentList: asttest.NewIdentList(declName), Type: tt.Expected}}
		sectionStr := "var " + declName + ": " + string(*tt.Text) + ";"
		sectionRunes := []rune(sectionStr)
		p := tt.newParser(&sectionRunes)
		res, err := p.ParseVarSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expectedSection, res)
		}
	})
	return tt
}
