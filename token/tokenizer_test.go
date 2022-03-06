package token_test

import (
	"testing"

	"github.com/akm/opparser/token"
)

func TestTokenizer(t *testing.T) {
	type pattern struct {
		text   string
		tokens TestTokens
	}
	patterns := []*pattern{
		{
			text: `unit Foo;
interface
type Bar = class // Bar is a class
end;
implementation
end.`,
			tokens: TestTokens{
				{Type: token.ReservedWord, Content: "unit"},
				{Type: token.Space, Content: " "},
				{Type: token.Identifier, Content: "Foo"},
				{Type: token.SpecialSymbol, Content: ";"},
				{Type: token.Space, Content: "\n"},
				{Type: token.ReservedWord, Content: "interface"},
				{Type: token.Space, Content: "\n"},
				{Type: token.ReservedWord, Content: "type"},
				{Type: token.Space, Content: " "},
				{Type: token.Identifier, Content: "Bar"},
				{Type: token.Space, Content: " "},
				{Type: token.SpecialSymbol, Content: "="},
				{Type: token.Space, Content: " "},
				{Type: token.ReservedWord, Content: "class"},
				{Type: token.Space, Content: " "},
				{Type: token.Comment, Content: "// Bar is a class"},
				{Type: token.Space, Content: "\n"},
				{Type: token.ReservedWord, Content: "end"},
				{Type: token.SpecialSymbol, Content: ";"},
				{Type: token.Space, Content: "\n"},
				{Type: token.ReservedWord, Content: "implementation"},
				{Type: token.Space, Content: "\n"},
				{Type: token.ReservedWord, Content: "end"},
				{Type: token.SpecialSymbol, Content: "."},
			},
		},
	}
	for _, ptn := range patterns {
		t.Run(ptn.text, func(t *testing.T) {
			check(t, ptn.text, &ptn.tokens)
		})
	}
}
