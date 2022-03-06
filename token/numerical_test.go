package token_test

import (
	"testing"

	"github.com/akm/opparser/token"
)

func TestNumerical(t *testing.T) {
	patterns := TestPatterns{
		{
			text:   `0`,
			tokens: TestTokens{{Type: token.NumeralInt, Content: "0"}},
		},
		{
			text:   `1`,
			tokens: TestTokens{{Type: token.NumeralInt, Content: "1"}},
		},
		{
			text:   `123`,
			tokens: TestTokens{{Type: token.NumeralInt, Content: "123"}},
		},
		{
			text:   `123.456`,
			tokens: TestTokens{{Type: token.NumeralReal, Content: "123.456"}},
		},
		{
			text:   `-2`,
			tokens: TestTokens{{Type: token.NumeralInt, Content: "-2"}},
		},
		{
			text:   `-9.9876`,
			tokens: TestTokens{{Type: token.NumeralReal, Content: "-9.9876"}},
		},
	}
	patterns.check(t)
}
