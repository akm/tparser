package parser

import (
	"testing"

	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestIdentClausee(t *testing.T) {
	text := []rune(`U1,U2,U3;`)

	parser := NewParser(&text)
	parser.NextToken()
	res, err := parser.ParseIdentList(';')
	if assert.NoError(t, err) {
		expect := asttest.NewIdentList(
			asttest.NewIdent("U1", asttest.NewIdentLocation(1, 1, 0, 3)),
			asttest.NewIdent("U2", asttest.NewIdentLocation(1, 4, 3, 6)),
			asttest.NewIdent("U3", asttest.NewIdentLocation(1, 7, 6, 9)),
		)
		assert.Equal(t, &expect, res)
	}
}
