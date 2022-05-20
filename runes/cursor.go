package runes

type Cursor struct {
	Text     *[]rune
	Len      int
	Position *Position
	wasLF    bool
}

const CursorEOF = rune(0)

func NewCursor(text *[]rune) *Cursor {
	return &Cursor{
		Text:     text,
		Len:      len(*text),
		Position: NewPosition(),
	}
}

func (c *Cursor) Clone() *Cursor {
	return &Cursor{
		Text:     c.Text,
		Len:      c.Len,
		Position: c.Position.Clone(),
	}
}

func (c *Cursor) Current() rune {
	return c.Seek(0)
}

func (c *Cursor) Seek(n int) rune {
	if c.Position.Index+n < c.Len {
		return (*c.Text)[c.Position.Index+n]
	}
	return CursorEOF
}

func (c *Cursor) Next() rune {
	if c.wasLF {
		c.Position.nextLine()
		c.wasLF = false
	}
	if c.Position.Index < c.Len {
		c.Position.inc()
	} else {
		return CursorEOF
	}
	c.Position.next()
	r := c.Seek(0)
	if r == '\n' {
		c.wasLF = true
	}
	return r
}
