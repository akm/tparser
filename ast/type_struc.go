package ast

import "github.com/akm/tparser/ast/astcore"

// - StrucType
//   ```
//   [PACKED] (ArrayType [PACKED]| SetType | FileType | RecType [PACKED])
//   ```
// NOTICE this is NOT "StructType" but "StrucType"
type StrucType interface {
	isStrucType()
	IsPacked() bool
	// inherits
	Type
}

// - ArrayType
//   ```
//   ARRAY ['[' OrdinalType ','... ']'] OF Type [PortabilityDirective]
//   ```
type ArrayType struct {
	// IndexTypes []OrdinalType
	IndexTypes []Type // OrdinalType or TypeId to OrdinalType
	BaseType   Type
	Packed     bool
	// implements
	StrucType
}

func (*ArrayType) isType()          {}
func (*ArrayType) isStrucType()     {}
func (m *ArrayType) IsPacked() bool { return m.Packed }
func (m *ArrayType) Children() Nodes {
	r := make(Nodes, len(m.IndexTypes)+1)
	for i, m := range m.IndexTypes {
		r[i] = m
	}
	r[len(m.IndexTypes)] = m.BaseType
	return r
}

// - SetType
//   ```
//   SET OF OrdinalType [PortabilityDirective]
//   ```
type SetType struct {
	OrdinalType
	Packed bool
	// implements
	StrucType
}

func (*SetType) isType()          {}
func (*SetType) isStrucType()     {}
func (m *SetType) IsPacked() bool { return m.Packed }
func (m *SetType) Children() Nodes {
	return Nodes{m.OrdinalType}
}

// - FileType
//   ```
//   FILE OF TypeId [PortabilityDirective]
//   ```
type FileType struct {
	*TypeId
	Packed bool
	// implements
	StrucType
}

func (*FileType) isType()          {}
func (*FileType) isStrucType()     {}
func (m *FileType) IsPacked() bool { return m.Packed }
func (m *FileType) Children() Nodes {
	return Nodes{m.TypeId}
}

// - RecType
//   ```
//   RECORD [FieldList] END [PortabilityDirective]
//   ```
type RecType struct {
	FieldList *FieldList
	Packed    bool
	// implements
	StrucType
}

func (*RecType) isType()          {}
func (*RecType) isStrucType()     {}
func (m *RecType) IsPacked() bool { return m.Packed }
func (m *RecType) Children() Nodes {
	return Nodes{m.FieldList}
}

// - FieldList
//   ```
//   FieldDecl ';'... [VariantSection] [';']
//   ```
type FieldList struct {
	FieldDecl      *FieldDecl
	VariantSection *VariantSection
	// implements
	Node
}

func (m *FieldList) Children() Nodes {
	return Nodes{m.FieldDecl, m.VariantSection}
}

// - FieldDecl
//   ```
//   IdentList ':' Type [PortabilityDirective]
//   ```
type FieldDecl struct {
	IdentList IdentList
	Type      Type
	//PortabilityDirective
	// implements
	astcore.DeclNode
}

func (m *FieldDecl) Children() Nodes {
	return Nodes{m.Type}
}

func (m *FieldDecl) ToDeclarations() astcore.Decls {
	return astcore.NewDeclarations(m.IdentList, m)
}

// - VariantSection
//   ```
//   CASE [Ident ':'] TypeId OF RecVariant ';'...
//   ```
type VariantSection struct {
	Ident       *Ident
	TypeId      *TypeId
	RecVariants RecVariants
	// implements
	Node
}

func (m *VariantSection) Children() Nodes {
	return Nodes{m.TypeId, m.RecVariants}
}

// implements Node
type RecVariants []*RecVariant

func (s RecVariants) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - RecVariant
//   ```
//   ConstExpr ','... ':' '(' [FieldList] ')'
//   ```
type RecVariant struct {
	ConstExprs ConstExprs
	FieldList  *FieldList
	// implements
	Node
}

func (m *RecVariant) Children() Nodes {
	return Nodes{m.ConstExprs, m.FieldList}
}
