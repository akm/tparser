package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeRange(t *testing.T) {
	t.Run("CodeRange", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &CodeBlockNode{})
	})
}
