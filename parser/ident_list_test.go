package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestIdentClausee(t *testing.T) {
	text := []rune(`U1,U2,U3;`)

	parser := NewParser(&text)
	res, err := parser.ParseIdentClause()
	if assert.NoError(t, err) {
		assert.Equal(t, &ast.UsesClause{
			"U1", "U2", "U3",
		}, res)
	}
}
