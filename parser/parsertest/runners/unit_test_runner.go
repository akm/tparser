package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunUnitTest(t *testing.T, name string, text []rune, expected *ast.Unit, args ...interface{}) {
	NewUnitTestRunner(t, name, text, expected, args...).Run()
}

type UnitTestRunner struct {
	*BaseTestRunner
	Expected          *ast.Unit
	ClearUnitDeclMaps bool
}

func NewUnitTestRunner(t *testing.T, name string, text []rune, expected *ast.Unit, args ...interface{}) *UnitTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	runnerFuncs, rest3 := FilterUnitTestRunnerFuncs(rest2)
	if len(rest3) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest3))
	}
	r := &UnitTestRunner{
		BaseTestRunner:    NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:          expected,
		ClearUnitDeclMaps: true,
	}
	runnerFuncs.Call(r)
	return r
}

func (tt *UnitTestRunner) Run() *UnitTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			args := tt.ParserArgFuncs.Results()
			p := NewTestUnitParser(tt.Text, args...)
			p.NextToken()
			r, err := p.ParseUnit()
			if r != nil {
				if tt.ClearUnitDeclMaps {
					asttest.ClearUnitDeclMaps(tt.T, r)
				}
			}
			return r, err
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertUnit(t, tt.Expected, actual.(*ast.Unit))
				}
			}
		},
	)
	return tt
}

type UnitTestRunnerFunc = func(*UnitTestRunner)
type UnitTestRunnerFuncs []UnitTestRunnerFunc

func FilterUnitTestRunnerFuncs(args []interface{}) (UnitTestRunnerFuncs, []interface{}) {
	r := UnitTestRunnerFuncs{}
	rest := []interface{}{}
	for _, arg := range args {
		if v, ok := arg.(UnitTestRunnerFunc); ok {
			r = append(r, v)
		} else {
			rest = append(rest, arg)
		}
	}
	return r, rest
}

func (s UnitTestRunnerFuncs) Call(tt *UnitTestRunner) {
	for _, f := range s {
		f(tt)
	}
}
