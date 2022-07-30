package ast

import (
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

// - Expression
//   ```
//   SimpleExpression [RelOp SimpleExpression]...
//   ```
type Expression struct {
	*SimpleExpression
	RelOpSimpleExpressions RelOpSimpleExpressions
}

var _ Node = (*Expression)(nil)

func NewExpression(arg interface{}) *Expression {
	switch v := arg.(type) {
	case *Expression:
		return v
	case *SimpleExpression:
		return &Expression{SimpleExpression: v}
	default:
		return &Expression{SimpleExpression: NewSimpleExpression(arg)}
	}
}

func (m *Expression) Children() Nodes {
	r := Nodes{m.SimpleExpression}
	if m.RelOpSimpleExpressions != nil {
		r = append(r, m.RelOpSimpleExpressions)
	}
	return r
}

// - ExprList
//   ```
//   Expression ','...
//   ```
type ExprList []*Expression

var _ Node = (ExprList)(nil)

func (s ExprList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

// - RelOp
//   ```
//   '>'
//   ```
//   ```
//   '<'
//   ```
//   ```
//   '<='
//   ```
//   ```
//   '>='
//   ```
//   ```
//   '='
//   ```
//   ```
//   '<>'
//   ```
//   ```
//   IN
//   ```
//   ```
//   IS
//   ```
type RelOpSimpleExpression struct {
	RelOp string // '>' | '<' | '<=' | '>=' | '=' | '<>' | "IN" | "IS" | "AS"
	*SimpleExpression
}

var _ Node = (*RelOpSimpleExpression)(nil)

func (m *RelOpSimpleExpression) Children() Nodes {
	return Nodes{m.SimpleExpression}
}

type RelOpSimpleExpressions []*RelOpSimpleExpression

var _ Node = (RelOpSimpleExpressions)(nil)

func (s RelOpSimpleExpressions) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

// - SimpleExpression
//   ```
//   ['+' | '-'] Term [AddOp Term]...
//   ```
type SimpleExpression struct {
	UnaryOp *string //  '+' | '-' or nil
	*Term
	AddOpTerms AddOpTerms
}

var _ Node = (*SimpleExpression)(nil)

func NewSimpleExpression(arg interface{}) *SimpleExpression {
	switch v := arg.(type) {
	case *SimpleExpression:
		return v
	case *Term:
		return &SimpleExpression{Term: v}
	default:
		return &SimpleExpression{Term: NewTerm(arg)}
	}
}

func (m *SimpleExpression) Children() Nodes {
	r := Nodes{m.Term}
	if m.AddOpTerms != nil {
		r = append(r, m.AddOpTerms)
	}
	return r
}

// - AddOp
//   ```
//   '+'
//   ```
//   ```
//   '-'
//   ```
//   ```
//   OR
//   ```
//   ```
//   XOR
//   ```
type AddOpTerm struct {
	AddOp string // '+' | '-' | "OR" | "XOR"
	*Term
}

var _ Node = (*AddOpTerm)(nil)

func (m *AddOpTerm) Children() Nodes { return Nodes{m.Term} }

type AddOpTerms []*AddOpTerm

var _ Node = (AddOpTerms)(nil)

func (s AddOpTerms) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

// - Term
//   ```
//   Factor [MulOp Factor]...
//   ```
type Term struct {
	Factor       Factor
	MulOpFactors MulOpFactors
}

var _ Node = (*Term)(nil)

func NewTerm(arg interface{}) *Term {
	switch v := arg.(type) {
	case *Term:
		return v
	case Factor:
		return &Term{Factor: v}
	default:
		return &Term{Factor: NewDesignatorFactor(arg)}
	}
}

func (m *Term) Children() Nodes {
	r := Nodes{m.Factor}
	if m.MulOpFactors != nil {
		r = append(r, m.MulOpFactors)
	}
	return r
}

// - MulOp
//   ```
//   '*'
//   ```
//   ```
//   '/'
//   ```
//   ```
//   DIV
//   ```
//   ```
//   MOD
//   ```
//   ```
//   AND
//   ```
//   ```
//   SHL
//   ```
//   ```
//   SHR
//   ```
//   ```
//   AS
//   ```
//
type MulOpFactor struct {
	MulOp  string // '*' | '/' | "DIV" | "MOD", "AND", "SHL", "SHR"
	Factor Factor
}

var _ Node = (*MulOpFactor)(nil)

func (m *MulOpFactor) Children() Nodes { return Nodes{m.Factor} }

type MulOpFactors []*MulOpFactor

var _ Node = (MulOpFactors)(nil)

func (s MulOpFactors) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

// - Factor
//   ```
//   Designator ['(' ExprList ')']
//   ```
//   ```
//   '@' Designator
//   ```
//   ```
//   Number
//   ```
//   ```
//   String
//   ```
//   ```
//   NIL
//   ```
//   ```
//   '(' Expression ')'
//   ```
//   ```
//   NOT Factor
//   ```
//   ```
//   SetConstructor
//   ```
//   ```
//   TypeId '(' Expression ')'
//   ```
type Factor interface {
	Node
	isFactor()
}

// Designator ['(' ExprList ')']

type DesignatorFactor struct {
	*Designator
	ExprList ExprList
}

var _ Factor = (*DesignatorFactor)(nil)

func NewDesignatorFactor(arg interface{}) *DesignatorFactor {
	switch v := arg.(type) {
	case *DesignatorFactor:
		return v
	default:
		return &DesignatorFactor{Designator: NewDesignator(arg)}
	}
}

func (*DesignatorFactor) isFactor() {}
func (m *DesignatorFactor) Children() Nodes {
	r := Nodes{m.Designator}
	if m.ExprList != nil {
		r = append(r, m.ExprList)
	}
	return r
}

// '@' Designator

type Address struct {
	*Designator
}

var _ Factor = (*Address)(nil)

func (m *Address) Children() Nodes { return Nodes{m.Designator} }
func (*Address) isFactor()         {}

// - Designator
//   ```
//   QualId ['.' Ident | '[' ExprList ']' | '^']...
//   ```
type Designator struct {
	*QualId
	Items DesignatorItems
}

var _ Node = (*Address)(nil)

func NewDesignator(arg interface{}) *Designator {
	switch v := arg.(type) {
	case *Designator:
		return v
	case *QualId:
		return &Designator{QualId: v}
	case Ident:
		return &Designator{QualId: NewQualId(nil, NewIdentRef(&v, nil))}
	case *Ident:
		return &Designator{QualId: NewQualId(nil, NewIdentRef(v, nil))}
	case IdentRef:
		return &Designator{QualId: NewQualId(nil, &v)}
	case *IdentRef:
		return &Designator{QualId: NewQualId(nil, v)}
	case token.Token:
		return NewDesignator(NewIdent(&v))
	case *token.Token:
		return NewDesignator(NewIdent(v))
	default:
		panic(errors.Errorf("Unsupported type %T for NewDesignator", arg))
	}
}

func (m *Designator) Children() Nodes {
	r := Nodes{m.QualId}
	if m.Items != nil {
		r = append(r, m.Items)
	}
	return r
}

type DesignatorItems []DesignatorItem

var _ Node = (DesignatorItems)(nil)

func (s DesignatorItems) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

type DesignatorItem interface {
	Node
	isDesignatorItem()
}

type DesignatorItemIdent struct {
	*Ident
}

var _ DesignatorItem = (*DesignatorItemIdent)(nil)

func NewDesignatorItemIdent(arg interface{}) *DesignatorItemIdent {
	switch v := arg.(type) {
	case *DesignatorItemIdent:
		return v
	case *Ident:
		return &DesignatorItemIdent{Ident: v}
	case *token.Token:
		return &DesignatorItemIdent{Ident: NewIdent(v)}
	default:
		panic(errors.Errorf("Unsupported type %T for NewDesignatorItemIdent", arg))
	}
}

func (*DesignatorItemIdent) isDesignatorItem() {}
func (m *DesignatorItemIdent) Children() Nodes {
	return Nodes{m.Ident}
}

type DesignatorItemExprList ExprList // Must implement DesignatorItem, and ancestor ExprList implements Node.

func (DesignatorItemExprList) isDesignatorItem() {}
func (s DesignatorItemExprList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

type DesignatorItemDereference struct {
}

var _ DesignatorItem = (*DesignatorItemDereference)(nil)

func (*DesignatorItemDereference) Children() Nodes   { return Nodes{} }
func (*DesignatorItemDereference) isDesignatorItem() {}

//   ```
//   Number
//   ```

type NumberFactor struct {
	Value string
}

var _ Factor = (*NumberFactor)(nil)

func NewNumber(v string) *NumberFactor { return &NumberFactor{Value: v} }
func (*NumberFactor) Children() Nodes  { return Nodes{} }
func (*NumberFactor) isFactor()        {}

//   ```
//   String
//   ```
type StringFactor struct {
	Value string
}

var _ Factor = (*StringFactor)(nil)

func NewString(v string) *StringFactor { return &StringFactor{Value: v} }
func (*StringFactor) Children() Nodes  { return Nodes{} }
func (*StringFactor) isFactor()        {}

// ValueFactor for true, false or other values

type ValueFactor struct {
	Value string
}

var _ Factor = (*ValueFactor)(nil)

func NewValueFactor(v string) *ValueFactor { return &ValueFactor{Value: v} }
func (*ValueFactor) Children() Nodes       { return Nodes{} }
func (*ValueFactor) isFactor()             {}

// Ninl

type Nil struct {
}

var _ Factor = (*Nil)(nil)

func NewNil() *Nil           { return &Nil{} }
func (*Nil) Children() Nodes { return Nodes{} }
func (*Nil) isFactor()       {}

// Parentheses
type Parentheses struct { // Round brackets
	Expression *Expression
}

var _ Factor = (*Parentheses)(nil)

func (m *Parentheses) Children() Nodes { return Nodes{m.Expression} }
func (*Parentheses) isFactor()         {}

//   ```
//   NOT Factor
//   ```

type Not struct {
	Factor
}

var _ Factor = (*Not)(nil)

func (m *Not) Children() Nodes { return Nodes{m.Factor} }
func (*Not) isFactor()         {}

// - SetConstructor
//   ```
//   '[' [SetElement ','...] ']'
//   ```
type SetConstructor struct {
	SetElements []*SetElement
}

var _ Factor = (*SetConstructor)(nil)

func (SetConstructor) isFactor() {}
func (m *SetConstructor) Children() Nodes {
	r := make(Nodes, len(m.SetElements))
	for i, v := range m.SetElements {
		r[i] = v
	}
	return r
}

// - SetElement
//   ```
//   Expression ['..' Expression]
//   ```
type SetElement struct {
	*Expression
	SubRangeEnd *Expression
}

var _ Node = (*SetElement)(nil)

func NewSetElement(expr *Expression) *SetElement {
	return &SetElement{Expression: expr}
}

func (m *SetElement) Children() Nodes {
	r := Nodes{m.Expression}
	if m.SubRangeEnd != nil {
		r = append(r, m.SubRangeEnd)
	}
	return r
}

//   ```
//   TypeId '(' Expression ')'
//   ```
type TypeCast struct {
	TypeId     *TypeId
	Expression *Expression
}

var _ Factor = (*TypeCast)(nil)

func (*TypeCast) isFactor()         {}
func (m *TypeCast) Children() Nodes { return Nodes{m.TypeId, m.Expression} }
