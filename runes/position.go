package runes

type Position struct {
	Index int
	Line  int
	Col   int
}

func NewPosition() *Position {
	return &Position{
		Index: 0,
		Line:  1,
		Col:   1,
	}
}

func (p *Position) inc() {
	p.Index++
}

func (p *Position) nextLine() {
	p.Line++
	p.Col = 1
}

func (p *Position) next() {
	p.Col++
}

func (p *Position) Clone() *Position {
	r := *p
	return &r
}
