package ext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsReverse(t *testing.T) {
	s := Strings{"a", "b", "c"}
	r := s.Reverse()
	assert.Equal(t, Strings{"c", "b", "a"}, r)
}
