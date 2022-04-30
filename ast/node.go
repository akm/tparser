package ast

type Node interface {
	Children() Nodes
}

type Nodes []Node

func (s Nodes) Children() Nodes {
	return s
}

type LeafNode struct {
}

func (*LeafNode) Children() Nodes {
	return nil
}
