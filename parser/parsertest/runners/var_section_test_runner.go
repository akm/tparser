package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

type VarSectionTestRunner struct {
	t        *testing.T
	name     string
	text     *[]rune
	expected ast.VarSection
	funcs    []func() interface{}
}

func NewVarSectionTestRunner(t *testing.T, name string, text []rune, expected ast.VarSection, funcs ...func() interface{}) *VarSectionTestRunner {
	return &VarSectionTestRunner{t: t, name: name, text: &text, expected: expected, funcs: funcs}
}

func (tt *VarSectionTestRunner) newParser() *parser.Parser {
	args := make([]interface{}, len(tt.funcs))
	for i, f := range tt.funcs {
		args[i] = f()
	}
	r := NewTestParser(tt.text, args...)
	r.NextToken()
	return r
}

func (tt *VarSectionTestRunner) Run() *VarSectionTestRunner {
	tt.t.Run(tt.name, func(t *testing.T) {
		p := tt.newParser()
		res, err := p.ParseVarSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			if !assert.Equal(t, tt.expected, res) {
				asttest.AssertVarSection(t, tt.expected, res)
			}
		}
	})
	return tt
}
