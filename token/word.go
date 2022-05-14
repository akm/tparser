package token

import (
	"github.com/akm/tparser/runes"
)

func ProcessWord(c *runes.Cursor) *Token {
	if runes.IsWordHead(c.Current()) {
		start := c.Position.Clone()
		for runes.IsWord(c.Next()) {
		}
		t := NewToken(Identifier, c.Text, start, c.Position.Clone())
		s := t.Value()
		if isReservedWord(s) {
			t.Type = ReservedWord
		} else {
			t.Type = Identifier
		}
		return t
	}
	return nil
}
