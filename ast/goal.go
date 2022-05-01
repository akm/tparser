package ast

// - Goal
//   ```
//   (Program | Package | Library | Unit)
//   ```
type Goal interface {
	isGoal()
	GetPath() string
}
