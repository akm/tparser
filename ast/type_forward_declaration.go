package ast

type ForwardDeclaration interface {
	SetActualType(Type) error
}
