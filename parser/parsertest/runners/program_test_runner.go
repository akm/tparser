package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunProgramTest(t *testing.T, name string, text []rune, expected *ast.Program, args ...interface{}) {
	NewProgramTestRunner(t, name, text, expected, args...).Run()
}

type ProgramTestRunner struct {
	*BaseTestRunner
	Expected *ast.Program
}

func NewProgramTestRunner(t *testing.T, name string, text []rune, expected *ast.Program, args ...interface{}) *ProgramTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	if len(rest2) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest2))
	}
	r := &ProgramTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
	return r
}

func (tt *ProgramTestRunner) Run() *ProgramTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			args := tt.ParserArgFuncs.Results()
			p := NewTestProgramParser(tt.Text, args...)
			p.NextToken()
			r, err := p.ParseProgram()
			return r, err
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertProgram(t, tt.Expected, actual.(*ast.Program))
				}
			}
		},
	)
	return tt
}
