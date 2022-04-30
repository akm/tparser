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

func WalkDown(n Node, f func(Node) error) error {
	if err := f(n); err != nil {
		return err
	}
	for _, m := range n.Children() {
		if err := WalkDown(m, f); err != nil {
			return err
		}
	}
	return nil
}
