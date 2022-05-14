package ast

import "github.com/akm/tparser/ast/astcore"

// - CompoundStmt
//   ```
//   BEGIN StmtList END
//   ```
type CompoundStmt struct {
	StructStmt
	StmtList
}

func (*CompoundStmt) isStructStmt()     {}
func (m *CompoundStmt) Children() Nodes { return Nodes{m.StmtList} }

// - StmtList
//   ```
//   (Statement ';') ...
//   ```
type StmtList []*Statement // must implement Node
func (s StmtList) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

// - Statement
//   ```
//   [LabelId ':'] [SimpleStatement | StructStmt]
//   ```
type Statement struct {
	Node
	LabelId *LabelId
	Body    StatementBody
}

func (m *Statement) Children() Nodes {
	res := Nodes{}
	if m.LabelId != nil {
		res = append(res, m.LabelId)
	}
	res = append(res, m.Body)
	return res
}

// - SimpleStatement
//   (CallStatement)
//   ```
//   Designator ['(' [ExprList] ')']
//   ```
//   (AssignStatement)
//   ```
//   Designator ':=' Expression
//   ```
//   (InheritedStatement)
//   ```
//   INHERITED
//   ```
//   (GotoStatement)
//   ```
//   GOTO LabelId
//   ```
// - StructStmt
//   ```
//   CompoundStmt
//   ```
//   ```
//   ConditionalStmt
//   ```
//   ```
//   LoopStmt
//   ```
//   ```
//   WithStmt
//   ```
//   ```
//   TryExceptStmt
//   ```
//   ```
//   TryFinallyStmt
//   ```
//   ```
//   RaiseStmt
//   ```
//   ```
//   AssemblerStmt
//   ```
type StatementBody interface {
	Node
	isStatementBody()
}

type SimpleStatement interface {
	StatementBody
	isSimpleStatement()
}

type StructStmt interface {
	StatementBody
	isStructStmt()
}

// DesignatorStatement is an interface for CallStatement and AssignStatement
type DesignatorStatement interface {
	SimpleStatement
	isDesignatorStatement()
}

//   (CallStatement)
//   ```
//   Designator ['(' [ExprList] ')']
//   ```
type CallStatement struct {
	DesignatorStatement
	Designator *Designator
	ExprList   ExprList // nil able
}

func (*CallStatement) isStatementBody()         {}
func (*CallStatement) isSimpleStatement()       {}
func (m *CallStatement) isDesignatorStatement() {}
func (m *CallStatement) Children() Nodes {
	res := Nodes{m.Designator}
	if m.ExprList != nil {
		res = append(res, m.ExprList)
	}
	return res
}

//   (AssignStatement)
//   ```
//   Designator ':=' Expression
//   ```
type AssignStatement struct {
	DesignatorStatement
	Designator *Designator
	Expression *Expression
}

func (*AssignStatement) isStatementBody()         {}
func (*AssignStatement) isSimpleStatement()       {}
func (m *AssignStatement) isDesignatorStatement() {}
func (m *AssignStatement) Children() Nodes {
	return Nodes{m.Designator, m.Expression}
}

//   (InheritedStatement)
//   ```
//   INHERITED
//   ```
type InheritedStatement struct {
	SimpleStatement
	Ref *astcore.Declaration // reference to the ancestor method
}

func (*InheritedStatement) isStatementBody()   {}
func (*InheritedStatement) isSimpleStatement() {}
func (*InheritedStatement) Children() Nodes    { return Nodes{} }

//   (GotoStatement)
//   ```
//   GOTO LabelId
//   ```
type GotoStatement struct {
	SimpleStatement
	LabelId *LabelId
	Ref     *astcore.Declaration
}

func (*GotoStatement) isStatementBody()   {}
func (*GotoStatement) isSimpleStatement() {}
func (m *GotoStatement) Children() Nodes  { return Nodes{m.LabelId} }

// - ConditionalStmt
//   ```
//   IfStmt
//   ```
//   ```
//   CaseStmt
//   ```
type ConditionalStmt interface {
	StructStmt
	isConditionalStmt()
}

// - IfStmt
//   ```
//   IF Expression THEN Statement [ELSE Statement]
//   ```

type IfStmt struct {
	Condition *Expression
	Then      *Statement
	Else      *Statement
	ConditionalStmt
}

func (*IfStmt) isStatementBody()   {}
func (*IfStmt) isStructStmt()      {}
func (*IfStmt) isConditionalStmt() {}
func (m *IfStmt) Children() Nodes {
	r := Nodes{m.Condition, m.Then}
	if m.Else != nil {
		r = append(r, m.Else)
	}
	return r
}

// - CaseStmt
//   ```
//   CASE Expression OF CaseSelector ';'... [ELSE StmtList] [';'] END
//   ```

type CaseStmt struct {
	Expression *Expression
	Selectors  CaseSelectors
	Else       StmtList
}

func (*CaseStmt) isStatementBody()   {}
func (*CaseStmt) isStructStmt()      {}
func (*CaseStmt) isConditionalStmt() {}
func (m *CaseStmt) Children() Nodes {
	r := Nodes{m.Expression, m.Selectors}
	if m.Else != nil {
		r = append(r, m.Else)
	}
	return r
}

type CaseSelectors []*CaseSelector // must implements Node
func (s CaseSelectors) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

// - CaseSelector
//   ```
//   CaseLabel ','... ':' Statement
//   ```
type CaseSelector struct {
	Labels    CaseLabels
	Statement *Statement
	Node
}

func (m *CaseSelector) Children() Nodes {
	return Nodes{m.Labels, m.Statement}
}

type CaseLabels []*CaseLabel // must implements Node
func (s CaseLabels) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

// - CaseLabel
//   ```
//   ConstExpr ['..' ConstExpr]
//   ```

type CaseLabel struct {
	ConstExpr      *ConstExpr
	ExtraConstExpr *ConstExpr
	Node
}

func NewCaseLabel(expr *ConstExpr, extras ...*ConstExpr) *CaseLabel {
	switch len(extras) {
	case 0:
		return &CaseLabel{ConstExpr: expr}
	case 1:
		return &CaseLabel{ConstExpr: expr, ExtraConstExpr: extras[0]}
	default:
		panic("too many extras for NewCaseLabel")
	}
}

func (m *CaseLabel) Children() Nodes {
	r := Nodes{m.ConstExpr}
	if m.ExtraConstExpr != nil {
		r = append(r, m.ExtraConstExpr)
	}
	return r
}

// - LoopStmt
//   ```
//   RepeatStmt
//   ```
//   ```
//   WhileStmt
//   ```
//   ```
//   ForStmt
//   ```
type LoopStmt interface {
	StructStmt
	isLoopStmt()
}

// - RepeatStmt
//   ```
//   REPEAT StmtList UNTIL Expression
//   ```
type RepeatStmt struct {
	StmtList  StmtList
	Condition *Expression
	LoopStmt
}

func (*RepeatStmt) isStatementBody() {}
func (*RepeatStmt) isStructStmt()    {}
func (*RepeatStmt) isLoopStmt()      {}
func (m *RepeatStmt) Children() Nodes {
	return Nodes{m.StmtList, m.Condition}
}

// - WhileStmt
//   ```
//   WHILE Expression DO Statement
//   ```
type WhileStmt struct {
	Condition *Expression
	Statement *Statement
	LoopStmt
}

func (*WhileStmt) isStatementBody() {}
func (*WhileStmt) isStructStmt()    {}
func (*WhileStmt) isLoopStmt()      {}
func (m *WhileStmt) Children() Nodes {
	return Nodes{m.Condition, m.Statement}
}

// - ForStmt
//   ```
//   FOR QualId ':=' Expression (TO | DOWNTO) Expression DO Statement
//   ```
type ForStmt struct {
	QualId    *QualId
	Initial   *Expression
	Terminal  *Expression
	Down      bool // false: TO, true: DOWNTO
	Statement *Statement
	LoopStmt
}

func (*ForStmt) isStatementBody() {}
func (*ForStmt) isStructStmt()    {}
func (*ForStmt) isLoopStmt()      {}
func (m *ForStmt) Children() Nodes {
	return Nodes{m.QualId, m.Initial, m.Terminal, m.Statement}
}

// - WithStmt
//   ```
//   WITH IdentList DO Statement
//   ```
//   Ident in IdentList doesn't have Ref to Declaration. So we use QualIds instead.
type WithStmt struct {
	Objects   QualIds
	Statement *Statement
	StructStmt
}

func (*WithStmt) isStatementBody() {}
func (*WithStmt) isStructStmt()    {}
func (m *WithStmt) Children() Nodes {
	return Nodes{m.Objects, m.Statement}
}

type TryStmt interface {
	StructStmt
	isTryStmt()
}

// - TryExceptStmt
//   ```
//   TRY
//     Statement...
//   EXCEPT
//     ExceptionBlock
//   END
//   ```
type TryExceptStmt struct {
	Statements     StmtList
	ExceptionBlock *ExceptionBlock
	TryStmt
}

func (*TryExceptStmt) isStatementBody() {}
func (*TryExceptStmt) isStructStmt()    {}
func (*TryExceptStmt) isTryStmt()       {}
func (m *TryExceptStmt) Children() Nodes {
	return Nodes{m.Statements, m.ExceptionBlock}
}

// - ExceptionBlock
//   ```
//   [ON [Ident ‘:’] TypeID DO Statement]...
//   [ELSE Statement...]
//   ```
type ExceptionBlock struct {
	Handlers ExceptionBlockHandlers
	Else     StmtList
	Node
}

func (m *ExceptionBlock) Children() Nodes {
	return Nodes{m.Handlers, m.Else}
}

type ExceptionBlockHandlers []*ExceptionBlockHandler // must implements Node
func (s ExceptionBlockHandlers) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

type ExceptionBlockHandler struct {
	Decl      *ExceptionBlockHandlerDecl
	Statement *Statement
	Node
}

func (m *ExceptionBlockHandler) Children() Nodes {
	return Nodes{m.Decl, m.Statement}
}

type ExceptionBlockHandlerDecl struct {
	Ident *Ident
	Type  Type
	astcore.Decl
}

func (m *ExceptionBlockHandlerDecl) Children() Nodes {
	r := Nodes{}
	if m.Ident != nil {
		r = append(r, m.Ident)
	}
	r = append(r, m.Type)
	return r
}

func (m *ExceptionBlockHandlerDecl) ToDeclarations() astcore.Declarations {
	return astcore.Declarations{astcore.NewDeclaration(m.Ident, m)}
}

// - TryFinallyStmt
//   ```
//   TRY
//     Statement
//   FINALLY
//     Statement
//   END
//   ```
type TryFinallyStmt struct {
	Statements1 StmtList
	Statements2 StmtList
	TryStmt
}

func (*TryFinallyStmt) isStatementBody() {}
func (*TryFinallyStmt) isStructStmt()    {}
func (*TryFinallyStmt) isTryStmt()       {}
func (m *TryFinallyStmt) Children() Nodes {
	return Nodes{m.Statements1, m.Statements2}
}

// - RaiseStmt
//   ```
//   RAISE [object] [AT address]
//   ```
type RaiseStmt struct {
	Object  *Expression
	Address *Expression
	StructStmt
}

func (*RaiseStmt) isStatementBody() {}
func (*RaiseStmt) isStructStmt()    {}
func (m *RaiseStmt) Children() Nodes {
	r := Nodes{}
	if m.Object != nil {
		r = append(r, m.Object)
	}
	if m.Address != nil {
		r = append(r, m.Address)
	}
	return r
}

// - AssemblerStatement
//   ```
//   ASM
//   <assemblylanguage>
//   END
//   ```
type AssemblerStatement struct {
	StructStmt
}

func (*AssemblerStatement) isStatementBody() {}
func (*AssemblerStatement) isStructStmt()    {}
func (*AssemblerStatement) Children() Nodes  { return Nodes{} }
