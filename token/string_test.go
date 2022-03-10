package token_test

import (
	"testing"

	"github.com/akm/tparser/token"
)

func TestString(t *testing.T) {
	patterns := TestPatterns{
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
	patterns.check(t)
}
