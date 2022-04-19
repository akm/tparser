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
	Ident                 Ident
	PortabilityDirective  *PortabilityDirective // optional
	InterfaceSection      *InterfaceSection
	ImplementationSection *ImplementationSection
	InitSection           *InitSection // optional
}

// - InterfaceSection
//   ```
//   INTERFACE
//   [UsesClause]
//   [InterfaceDecl]...
//   ```
type InterfaceSection struct {
	UsesClause     *UsesClause // optional
	InterfaceDecls []InterfaceDecl
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
	canBeInterfaceDecl()
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

// - UnitId
//   ```
//   <unit-identifier>
//   ```
type UnitId Ident

func NewUnitId(name interface{}) *UnitId {
	r := UnitId(NewIdent(name))
	return &r
}

func (u *UnitId) String() string {
	if u == nil {
		return ""
	} else {
		return string(*u)
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
