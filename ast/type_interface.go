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
	Heritage InterfaceHeritage
	Guid     *InterfaceGuid
	Members  InterfaceMemberList
	// implements
	InterfaceType
}

func (*CustomInterfaceType) isType()               {}
func (*CustomInterfaceType) IsInterfaceType() bool { return true }
func (m *CustomInterfaceType) Children() Nodes {
	r := Nodes{}
	if m.Heritage != nil {
		r = append(r, m.Heritage)
	}
	if m.Guid != nil {
		r = append(r, m.Guid)
	}
	if m.Members != nil {
		r = append(r, m.Members)
	}
	return r
}

// - InterfaceHeritage
//   ```
// '(' TypeId ',' ... ')'
//   ```
type InterfaceHeritage []*TypeId

func (s InterfaceHeritage) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - InterfaceGuid
//   ```
//   '[' ConstExpr of string ']'
//   ```
type InterfaceGuid struct {
	*ConstExpr
}

func (m *InterfaceGuid) Children() Nodes {
	return Nodes{m.ConstExpr}
}

// - InterfaceMemberList
//   ```
//   InterfaceMember ';'...
//   ```
type InterfaceMemberList []InterfaceMember

func (s InterfaceMemberList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - InterfaceMember
//   ```
//   InterfaceMethod
//   ```
//   ```
//   InterfaceProperty
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

func (m *InterfaceMethod) Children() Nodes {
	return Nodes{m.Heading}
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

func (m *InterfaceProperty) Children() Nodes {
	r := Nodes{m.Ident}
	if m.Interface != nil {
		r = append(r, m.Interface)
	}
	if m.Read != nil {
		r = append(r, m.Read)
	}
	if m.Write != nil {
		r = append(r, m.Write)
	}
	return r
}
