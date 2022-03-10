package token

import (
	"github.com/akm/tparser/runes"
)

func ProcessComment(c *runes.Cursor) *Token {
	switch c.Current() {
	case '{':
		start := c.Position.Clone()
		for {
			if r := c.Next(); r == '}' || r == runes.CursorEOF {
				break
			}
		}
		c.Next()
		return NewToken(Comment, c.Text, start, c.Position.Clone())
	case '/':
		if c.Seek(1) == '/' {
			start := c.Position.Clone()
			for {
				if r := c.Next(); r == '\n' || r == runes.CursorEOF {
					break
				}
			}
			// c.Next() // Don't include \n in comment
			return NewToken(Comment, c.Text, start, c.Position.Clone())
		}
		return nil
	case '(':
		if c.Seek(1) == '*' {
			start := c.Position.Clone()
			for {
				r := c.Next()
				if r == '*' && c.Seek(1) == ')' {
					c.Next()
					break
				}
				if r == runes.CursorEOF {
					break
				}
			}
			c.Next()
			return NewToken(Comment, c.Text, start, c.Position.Clone())
		}
		return nil
	}
	return nil
}
