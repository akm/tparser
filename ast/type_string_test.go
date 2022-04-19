package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringType(t *testing.T) {
	t.Run("StringType implements Type", func(t *testing.T) {
		var typ Type
		typ = &StringType{}
		assert.Implements(t, (*Type)(nil), typ)
	})
}
