package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunTypeTest(t *testing.T, name string, text []rune, expected ast.Type, args ...interface{}) {
	NewTypeTestRunner(t, name, text, expected, args...).Run()
}

type TypeTestRunner struct {
	*BaseTestRunner
	Expected ast.Type
}

func NewTypeTestRunner(t *testing.T, name string, text []rune, expected ast.Type, args ...interface{}) *TypeTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	runnerFuncs, rest3 := FilterTypeTestRunnerFuncs(rest2)
	if len(rest3) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest3))
	}
	r := &TypeTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
	for _, f := range runnerFuncs {
		f(r)
	}
	return r
}

func (tt *TypeTestRunner) Run() *TypeTestRunner {
	tt.BaseTestRunner.Run(
		func(p *parser.Parser) (astcore.Node, error) {
			return p.ParseType()
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertType(t, tt.Expected, actual.(ast.Type))
				}
			}
		},
	)
	return tt
}

func (tt *TypeTestRunner) RunTypeSection(declName string) *TypeTestRunner {
	sectionStr := "type " + declName + " = " + string(*tt.Text) + ";"
	sectionRunes := []rune(sectionStr)
	sectRunner := NewTypeSectionTestRunner(
		tt.T,
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
		tt.T,
		tt.Name+" in var section",
		sectionRunes,
		ast.VarSection{{IdentList: asttest.NewIdentList(declName), Type: tt.Expected}},
		tt.ParserArgFuncs...,
	)
	sectRunner.Run()
	return tt
}

type TypeTestRunnerFunc = func(*TypeTestRunner)
type TypeTestRunnerFuncs []TypeTestRunnerFunc

func FilterTypeTestRunnerFuncs(args []interface{}) (TypeTestRunnerFuncs, []interface{}) {
	r := TypeTestRunnerFuncs{}
	rest := []interface{}{}
	for _, arg := range args {
		if v, ok := arg.(TypeTestRunnerFunc); ok {
			r = append(r, v)
		} else {
			rest = append(rest, arg)
		}
	}
	return r, rest
}
