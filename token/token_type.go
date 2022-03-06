package token

type Type uint8

const (
	EOF Type = iota // 0
	Space
	Comment
	SpecialSymbol
	Identifier
	QualifiedIdentifier
	ReservedWord
	Directive
	NumeralInt
	NumeralReal
	Label
	CharacterString
)

var TypeNames = map[Type]string{
	EOF:                 "EOF",
	Space:               "space",
	Comment:             "comment",
	SpecialSymbol:       "special symbol",
	Identifier:          "identifier",
	QualifiedIdentifier: "qualified identifier",
	ReservedWord:        "reserved word",
	Directive:           "directive",
	NumeralInt:          "int",
	NumeralReal:         "real",
	Label:               "label",
	CharacterString:     "character string",
}

func (t Type) String() string {
	return TypeNames[t]
}
