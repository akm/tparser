package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

type TypeSectionTest struct {
	t        *testing.T
	name     string
	text     []rune
	expected ast.TypeSection
	funcs    []func() interface{}
}

func NewTypeSectionTest(t *testing.T, name string, text []rune, expected ast.TypeSection, funcs ...func() interface{}) *TypeSectionTest {
	return &TypeSectionTest{t: t, name: name, text: text, expected: expected, funcs: funcs}
}

func (tt *TypeSectionTest) newParser(text *[]rune) *parser.Parser {
	args := make([]interface{}, len(tt.funcs))
	for i, f := range tt.funcs {
		args[i] = f()
	}
	r := NewTestParser(&tt.text, args...)
	r.NextToken()
	return r
}

func (tt *TypeSectionTest) Run() *TypeSectionTest {
	tt.t.Run(tt.name, func(t *testing.T) {
		p := tt.newParser(&tt.text)
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
