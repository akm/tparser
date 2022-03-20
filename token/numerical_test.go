package token_test

import (
	"testing"

	"github.com/akm/tparser/token"
)

func TestNumerical(t *testing.T) {
	run := func(text string, tk *TestToken) {
		ptn := &TestPattern{text: text, tokens: TestTokens{tk}}
		ptn.check(t)
	}
	run(`0`, &TestToken{Type: token.NumeralInt, Content: "0"})
	run(`1`, &TestToken{Type: token.NumeralInt, Content: "1"})
	run(`123`, &TestToken{Type: token.NumeralInt, Content: "123"})
	run(`123.456`, &TestToken{Type: token.NumeralReal, Content: "123.456"})
	run(`-2`, &TestToken{Type: token.NumeralInt, Content: "-2"})
	run(`-9.9876`, &TestToken{Type: token.NumeralReal, Content: "-9.9876"})
}
