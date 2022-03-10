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

var TwoRunesSpecialSymbols = map[[2]rune]bool{
	{'(', '*'}: true,
	{'(', '.'}: true,
	{'*', ')'}: true,
	{'.', ')'}: true,
	{'.', '.'}: true,
	{'/', '/'}: true,
	{':', '='}: true,
	{'<', '='}: true,
	{'>', '='}: true,
	{'<', '>'}: true,
}

func ProcessSingleSpecialSymbol(c *runes.Cursor) *Token {
	if SingleSpecialSymbols[c.Current()] {
		start := c.Position.Clone()
		c.Next()
		return NewToken(SpecialSymbol, c.Text, start, c.Position.Clone())
	}
	return nil
}
