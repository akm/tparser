package token

import (
	"github.com/akm/opparser/ext"
)

var ReservedWords = ext.Strings{
	"and",
	"array",
	"as",
	"asm",
	"begin",
	"case",
	"class",
	"const",
	"constructor",
	"destructor",
	"dispinterface",
	"div",
	"do",
	"downto",
	"else",
	"end",
	"except",
	"exports",
	"file",
	"finalization",
	"finally",
	"for",
	"function",
	"goto",
	"if",
	"implementation",
	"in",
	"inherited",
	"initialization",
	"inline",
	"interface",
	"is",
	"label",
	"library3",
	"mod",
	"nil",
	"not",
	"object",
	"of",
	"or",
	"packed",
	"procedure",
	"program",
	"property",
	"raise",
	"record",
	"repeat",
	"resourcestring",
	"set",
	"shl",
	"shr",
	"string",
	"then",
	"threadvar",
	"to",
	"try",
	"type",
	"unit",
	"until",
	"uses",
	"var",
	"while",
	"with",
	"xor",
}.Set()