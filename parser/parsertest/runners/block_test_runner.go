package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunBlockTest(t *testing.T, name string, text []rune, expected *ast.Block, args ...interface{}) {
	NewBlockTestRunner(t, name, text, expected, args...).Run()
}

type BlockTestRunner struct {
	*BaseTestRunner
	Expected *ast.Block
}

func NewBlockTestRunner(t *testing.T, name string, text []rune, expected *ast.Block, args ...interface{}) *BlockTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	if len(rest2) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest2))
	}
	r := &BlockTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
	return r
}

func (tt *BlockTestRunner) Run() *BlockTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			args := tt.ParserArgFuncs.Results()
			p := NewTestProgramParser(tt.Text, args...)
			p.NextToken()
			r, err := p.ParseBlock()
			return r, err
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertBlock(t, tt.Expected, actual.(*ast.Block))
				}
			}
		},
	)
	return tt
}
