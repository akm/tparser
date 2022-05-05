package astcore

// interface which is implemented by all declaration types
type Decl interface {
	ToDeclarations() Declarations
}
