package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
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
	if len(rest2) > 0 {
		panic(errors.Errorf("unexpected arguments: %v", rest2))
	}
	return &TypeTestRunner{
		BaseTestRunner: NewBaseTestRunner(t, name, &text, true, parserArgFuncs, baseRunnerFuncs),
		Expected:       expected,
	}
}

func (tt *TypeTestRunner) Run() *TypeTestRunner {
	tt.BaseTestRunner.Run(
		func() (astcore.Node, error) {
			return tt.NewParser().ParseType()
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
		tt.ParserArgFuncs.Interfaces()...,
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
		tt.ParserArgFuncs.Interfaces()...,
	)
	sectRunner.Run()
	return tt
}
