package ast

type ObjectType interface {
	IsObjectType() bool
	// implements
	Type
}

// - ObjectType
//   ```
//   OBJECT [ObjHeritage] [ObjFieldList] [MethodList] END
//   ```
type CustomObjectType struct {
	TypeId *TypeId
	// implements
	PointerType
}

func (*CustomObjectType) isType()            {}
func (*CustomObjectType) IsObjectType() bool { return true }
func (m *CustomObjectType) Children() Nodes {
	return Nodes{m.TypeId}
}

// - ObjHeritage
//   ```
//   '(' QualId ')'
//   ```
type ObjHeritage struct {
	*QualId
}

func (m *ObjHeritage) Children() Nodes {
	return Nodes{m.QualId}
}

// - ObjFieldList
//   ```
//   (IdentList ':' Type) ';'
//   ```
type ObjFieldList struct {
	IdentList IdentList
	Type      Type
}

func (m *ObjFieldList) Children() Nodes {
	return Nodes{m.IdentList, m.Type}
}

type MethodDirective string

const (
	MdAbstract    MethodDirective = "ABSTRACT"
	MdVirtual     MethodDirective = "VIRTUAL"
	MdOverride    MethodDirective = "OVERRIDE"
	MdOverload    MethodDirective = "OVERLOAD"
	MdReintroduce MethodDirective = "REINTRODUCE"
)

type MethodDirectives []MethodDirective

// - MethodList
//   ```
//   (MethodHeading [';' VIRTUAL]) ';'...
//   ```
type MethodList []*Method

func (s MethodList) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

type Method struct {
	MethodHeading MethodHeading
	Directives    MethodDirectives
	// implements
	Node
}

func (m *Method) Children() Nodes {
	return Nodes{m.MethodHeading}
}

// - MethodHeading
// ```
// ProcedureHeading
// ```
// ```
// FunctionHeading
// ```
// ```
// ConstructorHeading
// ```
// ```
// DestructorHeading
// ```
type MethodHeading interface {
	isMethodHeading()
	// implements
	Node
}

// - ConstructorHeading
//   ```
//   CONSTRUCTOR Ident [FormalParameters]
//   ```
type ConstructorHeading struct {
	*Ident
	FormalParameters FormalParameters
	// implements
	MethodHeading
}

// - DestructorHeading
//   ```
//   DESTRUCTOR Ident [FormalParameters]
//   ```
type DestructorHeading struct {
	*Ident
	FormalParameters FormalParameters
	// implements
	MethodHeading
}
