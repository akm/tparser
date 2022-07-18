package ast

import (
	"strings"

	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/log"
)

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

// - ForwardDeclaredClassType
//   ```
//   CLASS
//   ```
type ForwardDeclaredClassType struct {
	// This will be set at the end of actual class type declaration.
	Actual *CustomClassType
}

func (*ForwardDeclaredClassType) isType()           {}
func (*ForwardDeclaredClassType) IsClassType() bool { return true }
func (m *ForwardDeclaredClassType) Children() Nodes {
	return Nodes{}
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
func (m *CustomClassType) FindMemberDecl(name string) *astcore.Decl {
	defer log.TraceMethod("CustomClassType.FindMemberDecl: " + name)()

	kw := strings.ToLower(name)
	for _, m := range m.Members {
		if m.ClassFieldList != nil {
			for _, f := range m.ClassFieldList {
				if r := f.IdentList.Find(kw); r != nil {
					return f.ToDeclarations().Find(kw)
				}
			}
		}
		if m.ClassMethodList != nil {
			for _, method := range m.ClassMethodList {
				if strings.ToLower(method.Heading.GetIdent().Name) == kw {
					return method.ToDeclarations().Find(kw)
				}
			}
		}
		if m.ClassPropertyList != nil {
			for _, prop := range m.ClassPropertyList {
				if strings.ToLower(prop.Ident.Name) == kw {
					return prop.ToDeclarations().Find(kw)
				}
			}
		}
	}

	if m.Heritage != nil && len(m.Heritage) > 0 {
		parent := m.Heritage[0]
		if parent.Ref != nil {
			if typeDecl, ok := parent.Ref.Node.(*TypeDecl); ok {
				if parentClass, ok := typeDecl.Type.(*CustomClassType); ok {
					return parentClass.FindMemberDecl(name)
				}
			}
		}
	}
	return nil
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
	CvDefault   ClassVisibility = "default" // implicitly public
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
	ClassMethod bool
	Heading     ClassMethodHeading
	Directives  ClassMethodDirectiveList
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

type ClassMethodDirectiveList []ClassMethodDirective

func (s ClassMethodDirectiveList) Include(w string) bool {
	kw := ClassMethodDirective(strings.ToUpper(w))
	for _, m := range s {
		if m == kw {
			return true
		}
	}
	return false
}

var ClassMethodDirectives = ClassMethodDirectiveList{
	CmdAbstract,
	CmdVirtual,
	CmdOverride,
	CmdOverload,
	CmdReintroduce,
}

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
	// See "Property overrides and redeclarations" in Object Pascal Language Guide
	Parent *ClassProperty
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
//   [FormalParameters] ':' TypeId
//   ```
type PropertyInterface struct {
	Parameters FormalParameters
	Type       *TypeId
}

func (m *PropertyInterface) Children() Nodes {
	return Nodes{m.Parameters, m.Type}
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
