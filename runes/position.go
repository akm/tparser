package runes

import "fmt"

type Position struct {
	Line  int
	Col   int
	Index int
}

func NewPosition() *Position {
	return &Position{
		Line:  1,
		Col:   1,
		Index: 0,
	}
}

func (p *Position) inc() {
	p.Index++
}

func (p *Position) nextLine() {
	p.Line++
	p.Col = 0
}

func (p *Position) next() {
	p.Col++
}

func (p *Position) Clone() *Position {
	r := *p
	return &r
}

func (p *Position) String() string {
	return fmt.Sprintf("%d,%d", p.Line, p.Col)
}
