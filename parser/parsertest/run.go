package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

type TypeTest struct {
	t        *testing.T
	name     string
	text     []rune
	expected ast.Type
	funcs    []func() interface{}
}

func NewTypeTest(t *testing.T, name string, text []rune, expected ast.Type, funcs ...func() interface{}) *TypeTest {
	return &TypeTest{t: t, name: name, text: text, expected: expected, funcs: funcs}
}

func (tt *TypeTest) Run() *TypeTest {
	tt.t.Run(tt.name, func(t *testing.T) {
		args := make([]interface{}, len(tt.funcs))
		for i, f := range tt.funcs {
			args[i] = f()
		}
		parser := NewTestParser(&tt.text, args...)
		parser.NextToken()
		res, err := parser.ParseType()
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, tt.expected, res)
		}
	})
	return tt
}

func (tt *TypeTest) RunTypeSection(declName string) *TypeTest {
	tt.t.Run(tt.name+" in type section", func(t *testing.T) {
		args := make([]interface{}, len(tt.funcs))
		for i, f := range tt.funcs {
			args[i] = f()
		}
		expectedSection := ast.TypeSection{{Ident: asttest.NewIdent(declName), Type: tt.expected}}
		sectionStr := "type " + declName + " = " + string(tt.text) + ";"
		sectionRunes := []rune(sectionStr)
		parser := NewTestParser(&sectionRunes, args...)
		parser.NextToken()
		res, err := parser.ParseTypeSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expectedSection, res)
		}
	})
	return tt
}

func (tt *TypeTest) RunVarSection(declName string) *TypeTest {
	tt.t.Run(tt.name+" in var section", func(t *testing.T) {
		args := make([]interface{}, len(tt.funcs))
		for i, f := range tt.funcs {
			args[i] = f()
		}
		expectedSection := ast.VarSection{{IdentList: asttest.NewIdentList(declName), Type: tt.expected}}
		sectionStr := "var " + declName + ": " + string(tt.text) + ";"
		sectionRunes := []rune(sectionStr)
		parser := NewTestParser(&sectionRunes, args...)
		parser.NextToken()
		res, err := parser.ParseVarSection(true)
		if assert.NoError(t, err) {
			asttest.ClearLocations(t, res)
			assert.Equal(t, expectedSection, res)
		}
	})
	return tt
}
