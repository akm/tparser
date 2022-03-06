package token_test

import (
	"testing"

	"github.com/akm/opparser/token"
	"github.com/stretchr/testify/assert"
)

type TestToken struct {
	Type    token.Type
	Content string
}

func toTestToken(t *token.Token) *TestToken {
	return &TestToken{
		Type:    t.Type,
		Content: string(t.Raw()),
	}
}

type TestTokens []*TestToken

func ToTestTokens(tokens *[]token.Token) *TestTokens {
	res := TestTokens{}
	for _, t := range *tokens {
		res = append(res, toTestToken(&t))
	}
	return &res
}

func tokennize(text string) *[]token.Token {
	code := []rune(text)
	res := []token.Token{}
	x := token.NewTokenizer(&code)
	for {
		t := x.Next()
		if t.Type == token.EOF {
			break
		}
		res = append(res, *t)
	}
	return &res
}

func check(t *testing.T, text string, expected *TestTokens) {
	t.Logf("Running test for %s\n", text)
	tokens := tokennize(text)
	assert.Equal(t, *expected, *ToTestTokens(tokens))
}
