package token_test

import (
	"testing"

	"github.com/akm/tparser/token"
)

func TestSpecialSymbolDots(t *testing.T) {
	patterns := TestPatterns{
		{
			text: `1.2.3`,
			tokens: TestTokens{
				{Type: token.NumeralReal, Content: "1.2"},
				{Type: token.SpecialSymbol, Content: "."},
				{Type: token.NumeralInt, Content: "3"},
			},
		},
		{
			text: `1.2..3.4`,
			tokens: TestTokens{
				{Type: token.NumeralReal, Content: "1.2"},
				{Type: token.SpecialSymbol, Content: ".."},
				{Type: token.NumeralReal, Content: "3.4"},
			},
		},
		{
			text: `'A'..'Z'`,
			tokens: TestTokens{
				{Type: token.CharacterString, Content: "'A'"},
				{Type: token.SpecialSymbol, Content: ".."},
				{Type: token.CharacterString, Content: "'Z'"},
			},
		},
		{
			text: `Green..White`,
			tokens: TestTokens{
				{Type: token.Identifier, Content: "Green"},
				{Type: token.SpecialSymbol, Content: ".."},
				{Type: token.Identifier, Content: "White"},
			},
		},
		{
			text: `S + T`,
			tokens: TestTokens{
				{Type: token.Identifier, Content: "S"},
				{Type: token.SpecialSymbol, Content: "+"},
				{Type: token.Identifier, Content: "T"},
			},
		},
		{
			text: `S - T`,
			tokens: TestTokens{
				{Type: token.Identifier, Content: "S"},
				{Type: token.SpecialSymbol, Content: "-"},
				{Type: token.Identifier, Content: "T"},
			},
		},
	}
	patterns.check(t)
}
