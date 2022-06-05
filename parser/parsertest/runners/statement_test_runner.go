package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunStatementTest(t *testing.T, name string, text []rune, expected *ast.Statement, args ...interface{}) {
	NewStatementTestRunner(t, name, text, expected, args...).Run()
}

type StatementTestRunner struct {
	*BaseTestRunner
	Expected *ast.Statement
}

func NewStatementTestRunner(t *testing.T, name string, text []rune, expected *ast.Statement, args ...interface{}) *StatementTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	if len(rest2) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest2))
	}
	return &StatementTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
}

func (tt *StatementTestRunner) Run() *StatementTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			return tt.NewParser().ParseStatement()
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertStatement(t, tt.Expected, actual.(*ast.Statement))
				}
			}
		},
	)
	return tt
}
