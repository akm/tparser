package token

import (
	"github.com/akm/opparser/runes"
)

func ProcessEof(c *runes.Cursor) *Token {
	if c.Current() == runes.CursorEOF {
		pos := c.Position.Clone()
		return NewToken(EOF, c.Text, pos, pos)
	}
	return nil
}
