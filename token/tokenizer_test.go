package token_test

import (
	"testing"

	"github.com/akm/tparser/token"
)

func TestTokenizer(t *testing.T) {
	unitText := `unit Foo;
interface
type Bar = class // Bar is a class
end;
implementation
end.`

	patterns := TestPatterns{
		{
			name:  "simple unit",
			flags: token.LoadSpace | token.LoadComment,
			text:  unitText,
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
		{
			name:  "simple unit without space",
			flags: token.LoadComment,
			text:  unitText,
			tokens: TestTokens{
				{Type: token.ReservedWord, Content: "unit"},
				{Type: token.Identifier, Content: "Foo"},
				{Type: token.SpecialSymbol, Content: ";"},
				{Type: token.ReservedWord, Content: "interface"},
				{Type: token.ReservedWord, Content: "type"},
				{Type: token.Identifier, Content: "Bar"},
				{Type: token.SpecialSymbol, Content: "="},
				{Type: token.ReservedWord, Content: "class"},
				{Type: token.Comment, Content: "// Bar is a class"},
				{Type: token.ReservedWord, Content: "end"},
				{Type: token.SpecialSymbol, Content: ";"},
				{Type: token.ReservedWord, Content: "implementation"},
				{Type: token.ReservedWord, Content: "end"},
				{Type: token.SpecialSymbol, Content: "."},
			},
		},
		{
			name: "simple unit without space and comment",
			text: unitText,
			tokens: TestTokens{
				{Type: token.ReservedWord, Content: "unit"},
				{Type: token.Identifier, Content: "Foo"},
				{Type: token.SpecialSymbol, Content: ";"},
				{Type: token.ReservedWord, Content: "interface"},
				{Type: token.ReservedWord, Content: "type"},
				{Type: token.Identifier, Content: "Bar"},
				{Type: token.SpecialSymbol, Content: "="},
				{Type: token.ReservedWord, Content: "class"},
				{Type: token.ReservedWord, Content: "end"},
				{Type: token.SpecialSymbol, Content: ";"},
				{Type: token.ReservedWord, Content: "implementation"},
				{Type: token.ReservedWord, Content: "end"},
				{Type: token.SpecialSymbol, Content: "."},
			},
		},
	}
	patterns.check(t)
}
