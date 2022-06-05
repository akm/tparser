package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunTypeTest(t *testing.T, name string, text []rune, expected ast.Type, args ...interface{}) {
	NewTypeTestRunner(t, name, text, expected, args...).Run()
}

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
			if !assert.Equal(t, tt.Expected, res) {
				asttest.AssertType(t, tt.Expected, res)
			}
		}
	})
	return tt
}

func (tt *TypeTestRunner) RunTypeSection(declName string) *TypeTestRunner {
	sectionStr := "type " + declName + " = " + string(*tt.Text) + ";"
	sectionRunes := []rune(sectionStr)
	sectRunner := NewTypeSectionTestRunner(
		tt.t,
		tt.Name+" in type section",
		sectionRunes,
		ast.TypeSection{{Ident: asttest.NewIdent(declName), Type: tt.Expected}},
		tt.ParserArgFuncs...,
	)
	sectRunner.Run()
	return tt
}

func (tt *TypeTestRunner) RunVarSection(declName string) *TypeTestRunner {
	sectionStr := "var " + declName + ": " + string(*tt.Text) + ";"
	sectionRunes := []rune(sectionStr)
	sectRunner := NewVarSectionTestRunner(
		tt.t,
		tt.Name+" in var section",
		sectionRunes,
		ast.VarSection{{IdentList: asttest.NewIdentList(declName), Type: tt.Expected}},
		tt.ParserArgFuncs...,
	)
	sectRunner.Run()
	return tt
}
