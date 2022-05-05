package astcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdent(t *testing.T) {
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &Ident{})
		assert.Implements(t, (*Node)(nil), IdentList{})
	})
}
