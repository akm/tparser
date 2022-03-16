package token

import "github.com/akm/tparser/runes"

var SingleSpecialSymbols = map[rune]bool{
	'#':  true,
	'$':  true,
	'&':  true,
	'\'': true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	',':  true,
	'-':  true,
	'.':  true,
	'/':  true,
	':':  true,
	';':  true,
	'<':  true,
	'=':  true,
	'>':  true,
	'@':  true,
	'[':  true,
	']':  true,
	'^':  true,
	'{':  true,
	'}':  true,
}

// !, ", %, ?, \, _, | and ~ are not special symbols.

func ProcessSingleSpecialSymbol(c *runes.Cursor) *Token {
	if SingleSpecialSymbols[c.Current()] {
		start := c.Position.Clone()
		c.Next()
		return NewToken(SpecialSymbol, c.Text, start, c.Position.Clone())
	}
	return nil
}

var TwoRunesSpecialSymbols = map[[2]rune]bool{
	// {'(', '*'}: true, // comment
	{'(', '.'}: true,
	// {'*', ')'}: true, // comment
	{'.', ')'}: true,
	{'.', '.'}: true,
	// {'/', '/'}: true, // comment
	{':', '='}: true,
	{'<', '='}: true,
	{'>', '='}: true,
	{'<', '>'}: true,
}

func ProcessDoubleSpecialSymbol(c *runes.Cursor) *Token {
	for runes := range TwoRunesSpecialSymbols {
		if c.Current() == runes[0] && c.Seek(1) == runes[1] {
			start := c.Position.Clone()
			c.Next()
			c.Next()
			return NewToken(SpecialSymbol, c.Text, start, c.Position.Clone())
		}
	}
	return nil
}
