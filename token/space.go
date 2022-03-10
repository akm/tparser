package token

import (
	"unicode"

	"github.com/akm/tparser/runes"
)

func ProcessSpace(c *runes.Cursor) *Token {
	if unicode.IsSpace(c.Current()) {
		start := c.Position.Clone()
		for unicode.IsSpace(c.Next()) {
		}
		return NewToken(Space, c.Text, start, c.Position.Clone())
	}
	return nil
}
