package ast

import "github.com/akm/tparser/runes"

type CodePosition = runes.Position

type CodeRange struct {
	Path  string
	Start CodePosition
	End   CodePosition
}
