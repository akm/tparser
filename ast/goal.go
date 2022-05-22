package ast

import "github.com/akm/tparser/ast/astcore"

// - Goal
//   ```
//   (Program | Package | Library | Unit)
//   ```
type Goal interface {
	astcore.DeclNode
	isGoal()
	GetPath() string
}
