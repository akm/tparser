package ast

type InterfaceType interface {
	IsInterfaceType() bool
	// implements
	Type
}

// - InterfaceType
//   ```
//   INTERFACE
//   [InterfaceHeritage]
//   [InterfaceGuid]
//   [InterfaceMemberList]
//   ...
//   END
//   ```
type CustomInterfaceType struct {
	InterfaceHeritage *InterfaceHeritage
	Guid              *InterfaceGuid
	Members           InterfaceMemberList
	// implements
	InterfaceType
}

// - InterfaceHeritage
//   ```
// '(' TypeId ',' ... ')'
//   ```
type InterfaceHeritage []*TypeId

// - InterfaceGuid
//   ```
//   '[' ConstExpr of string ']'
//   ```
type InterfaceGuid struct {
	*ConstExpr
}

// - InterfaceMemberList
//   ```
//   InterfaceMember ';'...
//   ```
type InterfaceMemberList []*InterfaceMember

// - InterfaceMember
//   ```
//   InterfaceMethod
//   ```
//   ```
//   ClassProperty
//   ```
type InterfaceMember interface {
	isInterfaceMember()
	// implements
	Node
}

// - InterfaceMethod
//   ```
//   InterfaceMethodHeading; [InterfaceMethodDirective ';'...];
//   ```
type InterfaceMethod struct {
	Heading    InterfaceMethodHeading
	Directives InterfaceMethodDirectives
}

// - InterfaceMethodHeading
//   ```
//   ProcedureHeading
//   ```
//   ```
//   FunctionHeading
//   ```
type InterfaceMethodHeading interface {
	isInterfaceMethodHeading()
	// implements
	Node
}

// - InterfaceMethodDirective
//   ```
//   stdcall
//   ```
type InterfaceMethodDirective string

const (
	ImdStdcall InterfaceMethodDirective = "STDCALL"
)

type InterfaceMethodDirectives []InterfaceMethodDirective

// - InterfaceProperty
//   ```
//   PROPERTY Ident PropertyInterface [READ Ident] [WRITE Ident]
//   ```
type InterfaceProperty struct {
	Ident     *Ident
	Interface *PropertyInterface
	Read      *IdentRef
	Write     *IdentRef
}
