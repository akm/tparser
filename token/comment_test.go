package token_test

import (
	"testing"

	"github.com/akm/tparser/token"
)

func TestComment(t *testing.T) {
	patterns := TestPatterns{
		{
			flags:  token.LoadComment,
			text:   `// comment`,
			tokens: TestTokens{{Type: token.Comment, Content: "// comment"}},
		},
		{
			flags:  token.LoadComment,
			text:   `{ comment }`,
			tokens: TestTokens{{Type: token.Comment, Content: "{ comment }"}},
		},
		{
			flags:  token.LoadComment,
			text:   `(* comment *)`,
			tokens: TestTokens{{Type: token.Comment, Content: "(* comment *)"}},
		},
	}
	patterns.check(t)
}
