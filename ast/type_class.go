package ast

type ClassType interface {
	IsClassType() bool
	// implements
	Type
}

// - ClassType
//   ```
//   CLASS [ClassHeritage]
//   [ClassVisibility]
//   [ClassFieldList]
//   [ClassMethodList]
//   [ClassPropertyList]
//   END
//   ```
type CustomClassType struct {
	ClassHeritage *ClassHeritage
	Fields        ClassFieldList
	Methods       *ClassMethodList
	Properties    ClassPropertyList
}

// - ClassHeritage
//   ```
//   '(' IdentList ')'
//   ```
type ClassHeritage struct {
	IdentList IdentList
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
//   (ClassVisibility ObjFieldList) ';'...
//   ```
type ClassFieldList struct {
	Visibility ClassVisibility
	*ObjFieldList
}

// - ClassMethodList
//   ```
//   (ClassVisibility MethodList) ';'...
//   ```
type ClassMethodList struct {
	Visibility ClassVisibility
	MethodList
}

// - ClassPropertyList
//   ```
//   (ClassVisibility PropertyList ';')...
//   ```
type ClassPropertyList []*ClassPropertyList

type ClassProperty struct {
	Visibility ClassVisibility
	*Property
}

// - Property
//   ```
//   PROPERTY Ident [PropertyInterface] [PropertySpecifiers] [PortabilityDirective]
//   ```
type Property struct {
	Ident                *Ident
	PropertyInterface    *PropertyInterface
	PropertySpecifiers   *PropertySpecifiers
	PortabilityDirective PortabilityDirective
}

// - PropertyInterface
//   ```
//   [PropertyParameterList] ':' TypeId
//   ```
type PropertyInterface struct {
	PropertyParameterList *PropertyParameterList
}

// - PropertyParameterList
//   ```
//   '[' (IdentList ':' TypeId) ';'... ']'
//   ```
type PropertyParameterList []*PropertyParameter

type PropertyParameter struct {
	IdentList IdentList
	TypeId    *TypeId
}

// - PropertySpecifiers
//   ```
//   [INDEX ConstExpr]
//   [READ Ident]
//   [WRITE Ident]
//   [STORED (Ident | Constant)]
//   [(DEFAULT ConstExpr) | NODEFAULT]
//   [IMPLEMENTS TypeId]
//   ```
type PropertySpecifiers struct {
	Index      *ConstExpr
	Read       *IdentRef
	Write      *IdentRef
	Stored     *PropertyStoredDirective
	Default    *PropertyDefaultDirective
	Implements *TypeId
}

type PropertyStoredDirective struct {
	IdentRef *IdentRef
	Constant *bool
}

type PropertyDefaultDirective struct {
	Value     *ConstExpr
	NoDefault *bool
}
