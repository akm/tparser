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
type ObjFieldList []*ObjField

func (s ObjFieldList) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

type ObjField struct {
	IdentList IdentList
	Type      Type
}

func (m *ObjField) Children() Nodes {
	return Nodes{m.IdentList, m.Type}
}

// - MethodList
//   ```
//   (MethodHeading [';' VIRTUAL]) ';'...
//   ```
type MethodList struct {
	Methods MethodHeadings
	Virtual bool
	// implements
	Node
}

func (m *MethodList) Children() Nodes {
	return Nodes{m.Methods}
}

type MethodHeadings []MethodHeading

func (s MethodHeadings) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
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
