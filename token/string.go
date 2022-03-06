package token

import (
	"github.com/akm/opparser/runes"
)

func ProcessString(c *runes.Cursor) *Token {
	switch c.Current() {
	case '\'':
		start := c.Position.Clone()
		var last rune
		for {
			r := c.Next()
			if last != '\\' && r == '\'' {
				break
			}
			last = r
		}
		c.Next()
		return NewToken(CharacterString, c.Text, start, c.Position.Clone())
	}
	return nil
}
