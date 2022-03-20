package token

import (
	"testing"

	"github.com/akm/tparser/runes"
	"github.com/stretchr/testify/assert"
)

func TestPredicatorl(t *testing.T) {
	text0 := []rune("S - T")
	t0 := &Token{
		Type:  SpecialSymbol,
		text:  &text0,
		Start: &runes.Position{Index: 2, Line: 1, Col: 3},
		End:   &runes.Position{Index: 3, Line: 1, Col: 4},
	}
	p0 := Some(
		Symbol('+'),
		Symbol('-'),
		ReservedWord.HasKeyword("OR"),
		ReservedWord.HasKeyword("XOR"),
	)
	assert.True(t, p0.Predicate(t0))
}
