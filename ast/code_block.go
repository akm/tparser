package ast

type CodeBlock interface {
	Node
	GetRange() *CodeRange
}

type CodeBlockNode struct {
	Range *CodeRange
}

func (n *CodeBlockNode) GetRange() *CodeRange {
	return n.Range
}

func (n *CodeBlockNode) Children() []Node {
	return nil
}
