package ast

import "github.com/akm/tparser/ast/astcore"

// - Goal
//   ```
//   (Program | Package | Library | Unit)
//   ```
type Goal interface {
	astcore.Decl
	isGoal()
	GetPath() string
}
