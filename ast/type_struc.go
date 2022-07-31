package ast

import (
	"github.com/akm/tparser/ast/astcore"
)

// - StrucType
//   ```
//   [PACKED] (ArrayType [PACKED]| SetType | FileType | RecType [PACKED])
//   ```
// NOTICE this is NOT "StructType" but "StrucType"
type StrucType interface {
	Type
	isStrucType()
	IsPacked() bool
}

// - ArrayType
//   ```
//   ARRAY ['[' OrdinalType ','... ']'] OF Type [PortabilityDirective]
//   ```
type ArrayType struct {
	IndexTypes []OrdinalType
	BaseType   Type
	Packed     bool
}

var _ StrucType = (*ArrayType)(nil)

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
}

var _ StrucType = (*SetType)(nil)

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
}

var _ StrucType = (*FileType)(nil)

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
}

var _ StrucType = (*RecType)(nil)

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
	FieldDecls     FieldDecls
	VariantSection *VariantSection
}

var _ Node = (*FieldList)(nil)

func (m *FieldList) Children() Nodes {
	r := Nodes{m.FieldDecls}
	if m.VariantSection != nil {
		r = append(r, m.VariantSection)
	}
	return r
}

type FieldDecls []*FieldDecl

var _ Node = (FieldDecls)(nil)

func (s FieldDecls) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - FieldDecl
//   ```
//   IdentList ':' Type [PortabilityDirective]
//   ```
type FieldDecl struct {
	IdentList IdentList
	Type      Type
	//PortabilityDirective
}

var _ astcore.DeclNode = (*FieldDecl)(nil)

func (m *FieldDecl) Children() Nodes {
	return Nodes{m.IdentList, m.Type}
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
	TypeId      OrdinalType
	RecVariants RecVariants
}

var _ Node = (*VariantSection)(nil)

func (m *VariantSection) Children() Nodes {
	r := Nodes{m.TypeId}
	if m.Ident != nil {
		r = append(r, m.Ident)
	}
	if m.RecVariants != nil {
		r = append(r, m.RecVariants)
	}
	return r
}

type RecVariants []*RecVariant

var _ Node = (RecVariants)(nil)

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
}

var _ Node = (*RecVariant)(nil)

func (m *RecVariant) Children() Nodes {
	return Nodes{m.ConstExprs, m.FieldList}
}
