package ast

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
	ClassHeritage *ClassHeritage
	Members       ClassMemberSections
	// implements
	ClassType
}

// - ObjectType
//   ```
//   OBJECT [ClassHeritage]
//   [ClassMemberSection]
//   END
//   ```
type CustomObjectType struct {
	ClassHeritage ClassHeritage
	Members       ClassMemberSections
	// implements
	ObjectType
}

// - ClassHeritage
//   ```
//   '(' TypeId ',' ... ')'
//   ```
type ClassHeritage []*TypeId

// - ClassMemberSections
//   ```
//   ClassMemberSection ...
//   ```
type ClassMemberSections []*ClassMemberSection

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

// - ClassField
//   ```
//   IdentList ':' Type
//   ```
type ClassField struct {
	IdentList IdentList
	Type      Type
}

// - ClassMethodList
//   ```
//   ClassMethod ';'...
//   ```
type ClassMethodList []*ClassMethod

// - ClassMethod
//   ```
//   [CLASS] MethodHeading [';' ClassMethodDirective ...]
//   ```
type ClassMethod struct {
	Static bool
	ClassMethodHeading
	Directives ClassMethodDirectives
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

// - DestructorHeading
//   ```
//   DESTRUCTOR Ident [FormalParameters]
//   ```
type DestructorHeading struct {
	*Ident
	FormalParameters FormalParameters
	// implements
	ClassMethodHeading
}

// - ClassPropertyList
//   ```
//   ClassProperty ';' ...
//   ```
type ClassPropertyList []*ClassProperty

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
	PropertyInterface    *PropertyInterface
	Index                *ConstExpr
	Read                 *IdentRef
	Write                *IdentRef
	Stored               *PropertyStoredSpecifier
	Default              *PropertyDefaultSpecifier
	Implements           *TypeId
	PortabilityDirective PortabilityDirective
}

// - PropertyInterface
//   ```
//   [PropertyParameterList] ':' TypeId
//   ```
type PropertyInterface struct {
	PropertyParameterList *PropertyParameterList
	Type                  *TypeId
}

// - PropertyParameterList
//   ```
//   '[' PropertyParameter ';'... ']'
//   ```
type PropertyParameterList []*PropertyParameter

// - PropertyParameter
//   ```
//   IdentList ':' TypeId
//   ```
type PropertyParameter struct {
	IdentList IdentList
	TypeId    *TypeId
}

type PropertyStoredSpecifier struct {
	IdentRef *IdentRef
	Constant *bool
}

type PropertyDefaultSpecifier struct {
	Value     *ConstExpr
	NoDefault *bool
}
