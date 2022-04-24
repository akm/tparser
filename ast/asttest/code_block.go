package asttest

import "github.com/akm/tparser/ast"

func NewCodeBlockNode(start, end *ast.CodePosition) *ast.CodeBlockNode {
	return &ast.CodeBlockNode{
		Range: &ast.CodeRange{Start: *start, End: *end},
	}
}

func CodePosition(index, line, col int) *ast.CodePosition {
	return &ast.CodePosition{Index: index, Line: line, Col: col}
}
