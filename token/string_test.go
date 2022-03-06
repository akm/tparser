package token_test

import (
	"testing"

	"github.com/akm/opparser/token"
)

func TestString(t *testing.T) {
	type pattern struct {
		text   string
		tokens TestTokens
	}
	patterns := []*pattern{
		{
			text:   `''`,
			tokens: TestTokens{{Type: token.CharacterString, Content: "''"}},
		},
		{
			text:   `'string1'`,
			tokens: TestTokens{{Type: token.CharacterString, Content: "'string1'"}},
		},
		{
			text:   `'with \'single quotes\''`,
			tokens: TestTokens{{Type: token.CharacterString, Content: "'with \\'single quotes\\''"}},
		},
	}
	for _, ptn := range patterns {
		t.Run(ptn.text, func(t *testing.T) {
			check(t, ptn.text, &ptn.tokens)
		})
	}
}
