package astcore

// interface which is implemented by all declaration types
type DeclNode interface {
	Node
	ToDeclarations() Declarations
}
