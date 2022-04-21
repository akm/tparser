package ast

import "github.com/akm/tparser/runes"

type CodePosition = runes.Position

type CodeRange struct {
	Path  string
	Start CodePosition
	End   CodePosition
}

type CodeBlock interface {
	GetRange() *CodeRange
	Children() []CodeBlock
}

type CodeBlockNode struct {
	Range *CodeRange
}

func (n *CodeBlockNode) GetRange() *CodeRange {
	return n.Range
}

func (n *CodeBlockNode) Children() []CodeBlock {
	return nil
}
