package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	t.Run("ConstSection", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &ConstSection{})
	})
	t.Run("ConstantDecl", func(t *testing.T) {
		assert.Implements(t, (*CodeBlock)(nil), &ConstantDecl{})
	})
}
