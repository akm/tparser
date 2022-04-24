package ast

type Node interface {
	Children() []Node
}

type RootNode interface {
	Node
	GetPath() string
}
