package ast

type CodePosition struct {
	Line int
	Col  int
}

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
