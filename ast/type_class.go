package ast

import "github.com/akm/tparser/ast/astcore"

type ClassType interface {
	IsClassType() bool
	// implements
	Type
}

type ObjectType interface {
	IsObjectType() bool
	// implements
	Type
}

// - ClassType
//   ```
//   CLASS [ClassHeritage]
//   [ClassMemberSection]
//   END
//   ```
type CustomClassType struct {
	Heritage ClassHeritage
	Members  ClassMemberSections
	// implements
	ClassType
}

func (*CustomClassType) isType()           {}
func (*CustomClassType) IsClassType() bool { return true }
func (m *CustomClassType) Children() Nodes {
	r := Nodes{}
	if m.Heritage != nil {
		r = append(r, m.Heritage)
	}
	if m.Members != nil {
		r = append(r, m.Members)
	}
	return r
}

// - ObjectType
//   ```
//   OBJECT [ClassHeritage]
//   [ClassMemberSection]
//   END
//   ```
type CustomObjectType struct {
	Heritage ClassHeritage
	Members  ClassMemberSections
	// implements
	ObjectType
}

func (*CustomObjectType) isType()            {}
func (*CustomObjectType) IsObjectType() bool { return true }
func (m *CustomObjectType) Children() Nodes {
	r := Nodes{}
	if m.Heritage != nil {
		r = append(r, m.Heritage)
	}
	if m.Members != nil {
		r = append(r, m.Members)
	}
	return r
}

// - ClassHeritage
//   ```
//   '(' TypeId ',' ... ')'
//   ```
type ClassHeritage []*TypeId

func (s ClassHeritage) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ClassMemberSections
//   ```
//   ClassMemberSection ...
//   ```
type ClassMemberSections []*ClassMemberSection

func (s ClassMemberSections) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ClassMemberSection
//   ```
//   ClassVisibility
//   [ClassFieldList]
//   [ClassMethodList]
//   [ClassPropertyList]
//   ```
type ClassMemberSection struct {
	Visibility        ClassVisibility
	ClassFieldList    ClassFieldList
	ClassMethodList   ClassMethodList
	ClassPropertyList ClassPropertyList
}

func (m *ClassMemberSection) Children() Nodes {
	r := Nodes{}
	if m.ClassFieldList != nil {
		r = append(r, m.ClassFieldList)
	}
	if m.ClassMethodList != nil {
		r = append(r, m.ClassMethodList)
	}
	if m.ClassPropertyList != nil {
		r = append(r, m.ClassPropertyList)
	}
	return r
}

// - ClassVisibility
//   ```
//   [PUBLIC | PROTECTED | PRIVATE | PUBLISHED]
//   ```
type ClassVisibility string

const (
	CvPrivate   ClassVisibility = "PRIVATE"
	CvProtected ClassVisibility = "PROTECTED"
	CvPblic     ClassVisibility = "PUBLIC"
	CvPublished ClassVisibility = "PUBLISHED"
)

// - ClassFieldList
//   ```
//   (ClassField) ';'...
//   ```
type ClassFieldList []*ClassField

func (s ClassFieldList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ClassField
//   ```
//   IdentList ':' Type
//   ```
type ClassField struct {
	IdentList IdentList
	Type      Type
	// implements
	astcore.DeclNode
}

func (m *ClassField) Children() Nodes {
	return Nodes{m.IdentList, m.Type}
}
func (m *ClassField) ToDeclarations() astcore.Decls {
	r := make(astcore.Decls, len(m.IdentList))
	for i, ident := range m.IdentList {
		r[i] = &astcore.Decl{Ident: ident, Node: m}
	}
	return r
}

// - ClassMethodList
//   ```
//   ClassMethod ';'...
//   ```
type ClassMethodList []*ClassMethod

func (s ClassMethodList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ClassMethod
//   ```
//   [CLASS] ClassMethodHeading [';' ClassMethodDirective ...]
//   ```
type ClassMethod struct {
	Static     bool
	Heading    ClassMethodHeading
	Directives ClassMethodDirectives
	// implements
	astcore.DeclNode
}

func (m *ClassMethod) ToDeclarations() astcore.Decls {
	return astcore.Decls{{Ident: m.Heading.GetIdent(), Node: m}}
}

func (m *ClassMethod) Children() Nodes {
	return Nodes{m.Heading}
}

// - ClassMethodHeading
//   ```
//   ProcedureHeading
//   ```
//   ```
//   FunctionHeading
//   ```
//   ```
//   ConstructorHeading
//   ```
//   ```
//   DestructorHeading
//   ```
type ClassMethodHeading interface {
	GetIdent() *Ident
	isClassMethodHeading()
	// implements
	Node
}

// - ClassMethodDirective
//   ```
//   ABSTRACT
//   ```
//   ```
//   VIRTUAL
//   ```
//   ```
//   OVERRIDE
//   ```
//   ```
//   OVERLOAD
//   ```
//   ```
//   REINTRODUCE
//   ```
type ClassMethodDirective string

const (
	CmdAbstract    ClassMethodDirective = "ABSTRACT"
	CmdVirtual     ClassMethodDirective = "VIRTUAL"
	CmdOverride    ClassMethodDirective = "OVERRIDE"
	CmdOverload    ClassMethodDirective = "OVERLOAD"
	CmdReintroduce ClassMethodDirective = "REINTRODUCE"
)

type ClassMethodDirectives []ClassMethodDirective

// - ConstructorHeading
//   ```
//   CONSTRUCTOR Ident [FormalParameters]
//   ```

type ConstructorHeading struct {
	*Ident
	FormalParameters FormalParameters
	// implements
	ClassMethodHeading
}

func (m *ConstructorHeading) GetIdent() *Ident { return m.Ident }
func (m *ConstructorHeading) Children() Nodes {
	r := Nodes{m.Ident}
	if m.FormalParameters != nil {
		r = append(r, m.FormalParameters)
	}
	return r
}

// - DestructorHeading
//   ```
//   DESTRUCTOR Ident
//   ```
type DestructorHeading struct {
	*Ident
	// implements
	ClassMethodHeading
}

func (m *DestructorHeading) GetIdent() *Ident { return m.Ident }
func (m *DestructorHeading) Children() Nodes {
	return Nodes{m.Ident}
}

// - ClassPropertyList
//   ```
//   ClassProperty ';' ...
//   ```
type ClassPropertyList []*ClassProperty

func (s ClassPropertyList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ClassProperty
//   ```
// 	 PROPERTY Ident
//   [PropertyInterface]
//   [INDEX ConstExpr]
//   [READ Ident]
//   [WRITE Ident]
//   [STORED (Ident | Constant)]
//   [(DEFAULT ConstExpr) | NODEFAULT]
//   [IMPLEMENTS TypeId]
//   [PortabilityDirective]
//   ```
type ClassProperty struct {
	Ident                *Ident
	Interface            *PropertyInterface
	Index                *ConstExpr
	Read                 *IdentRef
	Write                *IdentRef
	Stored               *PropertyStoredSpecifier
	Default              *PropertyDefaultSpecifier
	Implements           *TypeId
	PortabilityDirective PortabilityDirective
	// implements
	astcore.DeclNode
}

func (m *ClassProperty) ToDeclarations() astcore.Decls {
	return astcore.Decls{{Ident: m.Ident, Node: m}}
}
func (m *ClassProperty) Children() Nodes {
	r := Nodes{m.Ident}
	if m.Interface != nil {
		r = append(r, m.Interface)
	}
	if m.Index != nil {
		r = append(r, m.Index)
	}
	if m.Read != nil {
		r = append(r, m.Read)
	}
	if m.Write != nil {
		r = append(r, m.Write)
	}
	if m.Stored != nil {
		r = append(r, m.Stored)
	}
	if m.Default != nil {
		r = append(r, m.Default)
	}
	if m.Implements != nil {
		r = append(r, m.Implements)
	}
	return r
}

// - PropertyInterface
//   ```
//   [PropertyParameterList] ':' TypeId
//   ```
type PropertyInterface struct {
	PropertyParameterList *PropertyParameterList
	Type                  *TypeId
}

func (m *PropertyInterface) Children() Nodes {
	return Nodes{m.PropertyParameterList, m.Type}
}

// - PropertyParameterList
//   ```
//   '[' PropertyParameter ';'... ']'
//   ```
type PropertyParameterList []*PropertyParameter

func (s PropertyParameterList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - PropertyParameter
//   ```
//   IdentList ':' TypeId
//   ```
type PropertyParameter struct {
	IdentList IdentList
	TypeId    *TypeId
}

func (m *PropertyParameter) Children() Nodes {
	return Nodes{m.IdentList, m.TypeId}
}

type PropertyStoredSpecifier struct {
	IdentRef *IdentRef
	Constant *bool
}

func (m *PropertyStoredSpecifier) Children() Nodes {
	return Nodes{m.IdentRef}
}

type PropertyDefaultSpecifier struct {
	Value     *ConstExpr
	NoDefault *bool
}

func (m *PropertyDefaultSpecifier) Children() Nodes {
	return Nodes{m.Value}
}
