package asttest

import "github.com/akm/tparser/ast"

func ClearRange(node ast.Node) {
	if b, ok := node.(ast.CodeBlock); ok {
		b.ClearRange()
	}
}

func ClearAllRange(node ast.Node) {
	ClearRange(node)
	for _, child := range node.Children() {
		ClearAllRange(child)
	}
}
