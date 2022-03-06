package token_test

import (
	"testing"

	"github.com/akm/opparser/token"
)

func TestComment(t *testing.T) {
	type pattern struct {
		text   string
		tokens TestTokens
	}
	patterns := []*pattern{
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
	for _, ptn := range patterns {
		t.Run(ptn.text, func(t *testing.T) {
			check(t, ptn.text, &ptn.tokens)
		})
	}
}
