package token

import (
	"github.com/akm/tparser/runes"
)

func ProcessNumeral(c *runes.Cursor) *Token {
	if runes.IsDigit(c.Current()) || (runes.IsUnaryOp(c.Current()) && runes.IsDigit(c.Seek(1))) {
		start := c.Position.Clone()
		tokenType := NumeralInt
		for {
			r := c.Next()
			if runes.IsDigit(r) {
				// OK
			} else if r == runes.CursorEOF {
				break
			} else if r == '.' {
				if c.Seek(1) == '.' {
					break
				} else {
					if tokenType == NumeralReal {
						break
					}
					tokenType = NumeralReal
				}
			} else {
				break
			}
		}
		return NewToken(tokenType, c.Text, start, c.Position.Clone())
	}
	return nil
}
