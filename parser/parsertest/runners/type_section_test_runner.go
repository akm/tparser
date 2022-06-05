package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

type TypeSectionTestRunner struct {
	t        *testing.T
	name     string
	text     *[]rune
	expected ast.TypeSection
	funcs    []func() interface{}
}

func NewTypeSectionTestRunner(t *testing.T, name string, text []rune, expected ast.TypeSection, funcs ...func() interface{}) *TypeSectionTestRunner {
	return &TypeSectionTestRunner{t: t, name: name, text: &text, expected: expected, funcs: funcs}
}

func (tt *TypeSectionTestRunner) newParser() *parser.Parser {
	args := make([]interface{}, len(tt.funcs))
	for i, f := range tt.funcs {
		args[i] = f()
	}
	r := NewTestParser(tt.text, args...)
	r.NextToken()
	return r
}

func (tt *TypeSectionTestRunner) Run() *TypeSectionTestRunner {
	tt.t.Run(tt.name, func(t *testing.T) {
		p := tt.newParser()
		res, err := p.ParseTypeSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			if !assert.Equal(t, tt.expected, res) {
				asttest.AssertTypeSection(t, tt.expected, res)
			}
		}
	})
	return tt
}
