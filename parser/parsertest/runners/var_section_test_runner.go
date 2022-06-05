package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type VarSectionTestRunner struct {
	*BaseTestRunner
	Expected ast.VarSection
}

func NewVarSectionTestRunner(t *testing.T, name string, text []rune, expected ast.VarSection, args ...interface{}) *VarSectionTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	if len(rest2) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest2))
	}
	return &VarSectionTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
}

func (tt *VarSectionTestRunner) Run() *VarSectionTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			return tt.NewParser().ParseVarSection(true)
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertVarSection(t, tt.Expected, actual.(ast.VarSection))
				}
			}
		},
	)
	return tt
}
