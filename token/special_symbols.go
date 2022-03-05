package token

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
