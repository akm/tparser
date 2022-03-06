package token_test

import (
	"testing"

	"github.com/akm/opparser/token"
)

func TestComment(t *testing.T) {
	patterns := TestPatterns{
		{
			text:   `// comment`,
			tokens: TestTokens{{Type: token.Comment, Content: "// comment"}},
		},
		{
			text:   `{ comment }`,
			tokens: TestTokens{{Type: token.Comment, Content: "{ comment }"}},
		},
		{
			text:   `(* comment *)`,
			tokens: TestTokens{{Type: token.Comment, Content: "(* comment *)"}},
		},
	}
	patterns.check(t)
}
