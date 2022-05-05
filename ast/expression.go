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

type ExprList []*Expression

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

func (m *RelOpSimpleExpression) Children() Nodes {
	return Nodes{m.SimpleExpression}
}

type RelOpSimpleExpressions []*RelOpSimpleExpression

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

func (m *AddOpTerm) Children() Nodes {
	return Nodes{m.Term}
}

type AddOpTerms []*AddOpTerm

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
	Factor
	MulOpFactors MulOpFactors
}

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
	MulOp string // '*' | '/' | "DIV" | "MOD", "AND", "SHL", "SHR"
	Factor
}

func (m *MulOpFactor) Children() Nodes {
	return Nodes{m.Factor}
}

type MulOpFactors []*MulOpFactor

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
func (*DesignatorFactor) isFactor() {}

type DesignatorFactor struct {
	Factor
	*Designator
	ExprList ExprList
}

func NewDesignatorFactor(arg interface{}) *DesignatorFactor {
	switch v := arg.(type) {
	case *DesignatorFactor:
		return v
	default:
		return &DesignatorFactor{Designator: NewDesignator(arg)}
	}
}

func (m *DesignatorFactor) Children() Nodes {
	r := Nodes{m.Designator}
	if m.ExprList != nil {
		r = append(r, m.ExprList)
	}
	return r
}

// '@' Designator
func (*Address) isFactor() {}

type Address struct {
	Factor
	*Designator
}

func (m *Address) Children() Nodes {
	return Nodes{m.Designator}
}

// - Designator
//   ```
//   QualId ['.' Ident | '[' ExprList ']' | '^']...
//   ```
type Designator struct {
	*QualId
	Items DesignatorItems
}

func NewDesignator(arg interface{}) *Designator {
	switch v := arg.(type) {
	case *Designator:
		return v
	case *QualId:
		return &Designator{QualId: v}
	case Ident:
		return &Designator{QualId: &QualId{Ident: &v}}
	case *Ident:
		return &Designator{QualId: &QualId{Ident: v}}
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

func (DesignatorItemIdent) isDesignatorItem() {}

type DesignatorItemIdent Ident // Must implement DesignatorItem, and ancestor Ident implements Node.

func NewDesignatorItemIdent(v interface{}) *DesignatorItemIdent {
	r := DesignatorItemIdent(*NewIdentFrom(v))
	return &r
}

func (m *DesignatorItemIdent) Children() Nodes {
	return Nodes{}
}

func (DesignatorItemExprList) isDesignatorItem() {}

type DesignatorItemExprList ExprList // Must implement DesignatorItem, and ancestor ExprList implements Node.

func (s DesignatorItemExprList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (*DesignatorItemDereference) isDesignatorItem() {}

type DesignatorItemDereference struct {
	DesignatorItem
}

func (*DesignatorItemDereference) Children() Nodes {
	return Nodes{}
}

//   ```
//   Number
//   ```
func (*Number) isFactor() {}

type Number struct {
	Factor
	Value string
}

func NewNumber(v string) *Number { return &Number{Value: v} }
func (*Number) Children() Nodes {
	return Nodes{}
}

//   ```
//   String
//   ```
func (*String) isFactor() {}

type String struct {
	Factor
	Value string
}

func NewString(v string) *String { return &String{Value: v} }
func (*String) Children() Nodes {
	return Nodes{}
}

// Ninl
func (*Nil) isFactor() {}

type Nil struct {
	Factor
}

func NewNil() *Nil { return &Nil{} }

func (*Nil) Children() Nodes {
	return Nodes{}
}

// Parentheses
func (*Parentheses) isFactor() {}

type Parentheses struct { // Round brackets
	Factor
	Expression *Expression
}

func (m *Parentheses) Children() Nodes {
	return Nodes{m.Expression}
}

//   ```
//   NOT Factor
//   ```
func (*Not) isFactor() {}

type Not struct {
	Factor
}

func (m *Not) Children() Nodes {
	return Nodes{m.Factor}
}

func (SetConstructor) isFactor() {}

// - SetConstructor
//   ```
//   '[' [SetElement ','...] ']'
//   ```
type SetConstructor struct {
	Factor
	SetElements []*SetElement
}

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

func NewSetElement(expr *Expression) *SetElement {
	return &SetElement{
		Expression: expr,
	}
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
func (*TypeCast) isFactor() {}

type TypeCast struct {
	Factor
	TypeId     *TypeId
	Expression Expression
}

func (m *TypeCast) Children() Nodes {
	return Nodes{m.TypeId, &m.Expression}
}
