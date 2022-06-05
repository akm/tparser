package runners

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

type TypeTestRunner struct {
	t        *testing.T
	Name     string
	Text     *[]rune
	Expected ast.Type
	Funcs    []func() interface{}
}

func NewTypeTestRunner(t *testing.T, name string, text []rune, expected ast.Type, funcs ...func() interface{}) *TypeTestRunner {
	return &TypeTestRunner{t: t, Name: name, Text: &text, Expected: expected, Funcs: funcs}
}

func (tt *TypeTestRunner) newParser(text *[]rune) *parser.Parser {
	args := make([]interface{}, len(tt.Funcs))
	for i, f := range tt.Funcs {
		args[i] = f()
	}
	r := NewTestParser(text, args...)
	r.NextToken()
	return r
}

func (tt *TypeTestRunner) Run() *TypeTestRunner {
	tt.t.Run(tt.Name, func(t *testing.T) {
		p := tt.newParser(tt.Text)
		res, err := p.ParseType()
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, tt.Expected, res)
		}
	})
	return tt
}

func (tt *TypeTestRunner) RunTypeSection(declName string) *TypeTestRunner {
	tt.t.Run(tt.Name+" in type section", func(t *testing.T) {
		expectedSection := ast.TypeSection{{Ident: asttest.NewIdent(declName), Type: tt.Expected}}
		sectionStr := "type " + declName + " = " + string(*tt.Text) + ";"
		sectionRunes := []rune(sectionStr)
		p := tt.newParser(&sectionRunes)
		res, err := p.ParseTypeSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expectedSection, res)
		}
	})
	return tt
}

func (tt *TypeTestRunner) RunVarSection(declName string) *TypeTestRunner {
	tt.t.Run(tt.Name+" in var section", func(t *testing.T) {
		expectedSection := ast.VarSection{{IdentList: asttest.NewIdentList(declName), Type: tt.Expected}}
		sectionStr := "var " + declName + ": " + string(*tt.Text) + ";"
		sectionRunes := []rune(sectionStr)
		p := tt.newParser(&sectionRunes)
		res, err := p.ParseVarSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expectedSection, res)
		}
	})
	return tt
}
