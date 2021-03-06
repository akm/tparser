package cursor_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/akm/tparser/runes"
	"github.com/stretchr/testify/assert"
)

func TestCursorLFTest(t *testing.T) {
	fp, err := os.Open("./cursor_lf.txt")
	assert.NoError(t, err)
	defer fp.Close()

	str, err := ioutil.ReadAll(fp)
	assert.NoError(t, err)
	source := []rune(string(str))

	assertAndNext := func(t *testing.T, c *runes.Cursor, expected rune, line, col, index int) {
		assert.Equal(t, expected, c.Current())
		assert.Equal(t, &runes.Position{Line: line, Col: col, Index: index}, c.Position)
		c.Next()
	}

	c := runes.NewCursor(&source)
	assertAndNext(t, c, 'f', 1, 1, 0)
	assertAndNext(t, c, 'o', 1, 2, 1)
	assertAndNext(t, c, 'o', 1, 3, 2)
	assertAndNext(t, c, '\n', 1, 4, 3)
	assertAndNext(t, c, 'b', 2, 1, 4)
	assertAndNext(t, c, 'a', 2, 2, 5)
	assertAndNext(t, c, 'r', 2, 3, 6)
	assertAndNext(t, c, '\n', 2, 4, 7)
	assertAndNext(t, c, '\n', 3, 1, 8)
	assertAndNext(t, c, '\n', 4, 1, 9)
	assertAndNext(t, c, 'b', 5, 1, 10)
	assertAndNext(t, c, 'a', 5, 2, 11)
	assertAndNext(t, c, 'z', 5, 3, 12)
	assertAndNext(t, c, '\n', 5, 4, 13)
	assertAndNext(t, c, runes.CursorEOF, 6, 1, 14)
	assertAndNext(t, c, runes.CursorEOF, 6, 1, 14)
}
