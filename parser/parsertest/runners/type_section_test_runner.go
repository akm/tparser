package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func RunTypeSection(t *testing.T, name string, text []rune, expected ast.TypeSection, args ...interface{}) {
	NewTypeSectionTestRunner(t, name, text, expected, args...).Run()
}

type TypeSectionTestRunner struct {
	*BaseTestRunner
	Expected ast.TypeSection
}

func NewTypeSectionTestRunner(t *testing.T, name string, text []rune, expected ast.TypeSection, args ...interface{}) *TypeSectionTestRunner {
	parserArgFuncs, rest1 := FilterParserArgFuncs(args)
	baseRunnerFuncs, rest2 := FilterBaseTestRunnerFuncs(rest1)
	if len(rest2) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest2))
	}
	return &TypeSectionTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
}

func (tt *TypeSectionTestRunner) Run() *TypeSectionTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			return tt.NewParser().ParseTypeSection(true)
		},
		func(t *testing.T, actual astcore.Node) {
			if !assert.Equal(t, tt.Expected, actual) {
				if assert.IsType(t, tt.Expected, actual) {
					asttest.AssertTypeSection(t, tt.Expected, actual.(ast.TypeSection))
				}
			}
		},
	)
	return tt
}
