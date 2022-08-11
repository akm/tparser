package astcore

type Node interface {
	Children() Nodes
}

type Nodes []Node

var _ Node = (Nodes)(nil)

func (s Nodes) Children() Nodes {
	return s
}

// This method doesn't work well and fails to distinguish whether i is really nil
// Becasue i typed Node is nil with original Type.
// func (s Nodes) Compact() Nodes {
// 	if s == nil {
// 		return Nodes{}
// 	}
// 	r := Nodes{}
// 	for _, i := range s {
// 		if i != nil {
// 			r = append(r, i)
// 		}
// 	}
// 	return r
// }

func WalkDown(n Node, f func(Node) error) error {
	if err := f(n); err != nil {
		return err
	}
	for _, m := range n.Children() {
		if m != nil {
			if err := WalkDown(m, f); err != nil {
				return err
			}
		}
	}
	return nil
}
