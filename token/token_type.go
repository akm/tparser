package token

import (
	"fmt"

	"github.com/pkg/errors"
)

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
	PortabilityDirective
	NumeralInt
	NumeralReal
	Label
	CharacterString
)

var TypeNames = map[Type]string{
	EOF:                  "EOF",
	Space:                "space",
	Comment:              "comment",
	SpecialSymbol:        "special symbol",
	Identifier:           "identifier",
	QualifiedIdentifier:  "qualified identifier",
	ReservedWord:         "reserved word",
	Directive:            "directive",
	PortabilityDirective: "portability directive",
	NumeralInt:           "int",
	NumeralReal:          "real",
	Label:                "label",
	CharacterString:      "character string",
}

func (t Type) String() string {
	return TypeNames[t]
}

// As Predicator

func (typ Type) Name() string {
	return TypeNames[typ]
}

func (typ Type) Predicate(t *Token) bool {
	return t.Type == typ
}

func (typ Type) HasText(s string) Predicator {
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("%s has %q", typ.String(), s),
		predicate: func(t *Token) bool { return t.Type == typ && t.Value() == s },
	}
}

// kw must be upper case
func (typ Type) HasKeyword(kw string) Predicator {
	if !isReservedWord(kw) {
		panic(errors.Errorf("%q is not a reserved word", kw))
	}
	return &TokenPredicateImpl{
		name:      fmt.Sprintf("%s has %q", typ.String(), kw),
		predicate: func(t *Token) bool { return t.Type == typ && t.Value() == kw },
	}
}
