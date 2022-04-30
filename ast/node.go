package ast

type Node interface {
	Children() Nodes
}

type Nodes []Node

func (s Nodes) Children() Nodes {
	return s
}

func (s Nodes) Compact() Nodes {
	r := Nodes{}
	for _, m := range s {
		if m != nil {
			r = append(r, m)
		}
	}
	return r
}

type LeafNode struct {
}

func (*LeafNode) Children() Nodes {
	return nil
}
