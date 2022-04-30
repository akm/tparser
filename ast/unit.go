package ast

func (*Unit) isGoal() {}

// - Unit
//   ```
//   UNIT Ident [PortabilityDirective] ';'
//   InterfaceSection
//   ImplementationSection
//   InitSection '.'
//   ```
type Unit struct {
	Path                  string
	Ident                 Ident
	PortabilityDirective  *PortabilityDirective // optional
	InterfaceSection      *InterfaceSection
	ImplementationSection *ImplementationSection
	InitSection           *InitSection // optional
}

func (m *Unit) GetPath() string {
	return m.Path
}

func (m *Unit) Children() Nodes {
	return Nodes{m.Ident, m.InterfaceSection, m.ImplementationSection, m.InitSection}.Compact()
}

// - InterfaceSection
//   ```
//   INTERFACE
//   [UsesClause]
//   [InterfaceDecl]...
//   ```
type InterfaceSection struct {
	UsesClause     *UsesClause // optional
	InterfaceDecls InterfaceDecls
}

func (m *InterfaceSection) Children() Nodes {
	return Nodes{m.UsesClause, m.InterfaceDecls}.Compact()
}

// - InterfaceDecl
//   ```
//   ConstSection
//   ```
//   ```
//   TypeSection
//   ```
//   ```
//   VarSection
//   ```
//   ```
//   ExportedHeading
//   ```
type InterfaceDecl interface {
	Node
	canBeInterfaceDecl()
}

type InterfaceDecls []InterfaceDecl

func (s InterfaceDecls) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ImplementationSection
//   ```
//   IMPLEMENTATION
//   [UsesClause]
//   [DeclSection]...
//   [ExportsStmt]...
//   ```
type ImplementationSection struct {
}

func (m *ImplementationSection) Children() Nodes {
	return Nodes{}
}

// - InitSection
//   ```
//   INITIALIZATION StmtList [FINALIZATION StmtList] END
//   ```
//   ```
//   BEGIN StmtList END
//   ```
//   ```
//   END
//   ```
type InitSection struct {
}

func (m *InitSection) Children() Nodes {
	return Nodes{}
}

// - UnitId
//   ```
//   <unit-identifier>
//   ```
type UnitId Ident

func NewUnitId(name interface{}) *UnitId {
	r := UnitId(*NewIdentFrom(name))
	return &r
}

func (u *UnitId) String() string {
	if u == nil {
		return ""
	} else {
		return u.Name
	}
}

// - QualId
//   ```
//   [UnitId '.'] Ident
//   ```
type QualId struct {
	UnitId *UnitId
	Ident  Ident
}

func (m *QualId) Children() Nodes {
	return Nodes{m.UnitId, m.Ident}.Compact()
}
