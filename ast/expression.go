package ast

import (
	"fmt"
	"strings"

	"github.com/akm/tparser/ext"
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
	fmt.Stringer
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

func (m *Expression) String() string {
	return fmt.Sprintf("%s %s", m.SimpleExpression.String(), m.RelOpSimpleExpressions.String())
}

// - ExprList
//   ```
//   Expression ','...
//   ```
type ExprList []*Expression // must implement Node, fmt.Stringer

func (s ExprList) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (s ExprList) String() string {
	r := make(ext.Strings, len(s))
	for idx, i := range s {
		r[idx] = i.String()
	}
	return strings.Join(r, ", ")
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
	fmt.Stringer
}

func (m *RelOpSimpleExpression) Children() Nodes {
	return Nodes{m.SimpleExpression}
}

func (m *RelOpSimpleExpression) String() string {
	return fmt.Sprintf("%s %s", m.RelOp, m.SimpleExpression.String())
}

type RelOpSimpleExpressions []*RelOpSimpleExpression // must implement Node, fmt.Stringer

func (s RelOpSimpleExpressions) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (s RelOpSimpleExpressions) String() string {
	r := make(ext.Strings, len(s))
	for idx, i := range s {
		r[idx] = i.String()
	}
	return strings.Join(r, " ")
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

func (m *SimpleExpression) String() string {
	return fmt.Sprintf("%s %s", m.AddOpTerms.String(), m.AddOpTerms.String())
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
	fmt.Stringer
}

func (m *AddOpTerm) Children() Nodes { return Nodes{m.Term} }

func (m *AddOpTerm) String() string {
	return fmt.Sprintf("%s %s", m.AddOp, m.Term.String())
}

type AddOpTerms []*AddOpTerm // must implement Node, fmt.Stringer

func (s AddOpTerms) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (s AddOpTerms) String() string {
	r := make(ext.Strings, len(s))
	for idx, i := range s {
		r[idx] = i.String()
	}
	return strings.Join(r, " ")
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

func (m *Term) String() string {
	return fmt.Sprintf("%s %s", m.Factor.String(), m.MulOpFactors.String())
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

func (m *MulOpFactor) Children() Nodes { return Nodes{m.Factor} }

func (m *MulOpFactor) String() string {
	return fmt.Sprintf("%s %s", m.MulOp, m.Factor.String())
}

type MulOpFactors []*MulOpFactor // must implement Node, fmt.Stringer

func (s MulOpFactors) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (s MulOpFactors) String() string {
	r := make(ext.Strings, len(s))
	for idx, i := range s {
		r[idx] = i.String()
	}
	return strings.Join(r, " ")
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
	isFactor()
	Node
	fmt.Stringer
}

// Designator ['(' ExprList ')']

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

func (*DesignatorFactor) isFactor() {}
func (m *DesignatorFactor) Children() Nodes {
	r := Nodes{m.Designator}
	if m.ExprList != nil {
		r = append(r, m.ExprList)
	}
	return r
}

func (m *DesignatorFactor) String() string {
	return fmt.Sprintf("%s(%s)", m.Factor.String(), m.ExprList.String())
}

// '@' Designator

type Address struct {
	Factor
	*Designator
}

func (m *Address) Children() Nodes { return Nodes{m.Designator} }
func (*Address) isFactor()         {}
func (m *Address) String() string {
	return fmt.Sprintf("@%s", m.Designator.String())
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
		return &Designator{QualId: NewQualId(nil, &v)}
	case *Ident:
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

func (m *Designator) String() string {
	return fmt.Sprintf("%s%s", m.QualId.String(), m.Items.String())
}

type DesignatorItems []DesignatorItem

func (s DesignatorItems) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (s DesignatorItems) String() string {
	r := make(ext.Strings, len(s))
	for idx, i := range s {
		r[idx] = i.String()
	}
	return strings.Join(r, ", ")
}

type DesignatorItem interface {
	isDesignatorItem()
	Node
	fmt.Stringer
}

type DesignatorItemIdent Ident // Must implement DesignatorItem, and ancestor Ident implements Node.

func NewDesignatorItemIdent(v interface{}) *DesignatorItemIdent {
	r := DesignatorItemIdent(*NewIdentFrom(v))
	return &r
}

func (m *DesignatorItemIdent) Children() Nodes { return Nodes{} }
func (DesignatorItemIdent) isDesignatorItem()  {}
func (m *DesignatorItemIdent) String() string {
	return fmt.Sprintf(".%s", m.Name)
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
func (m DesignatorItemExprList) String() string {
	return fmt.Sprintf("[%s]", m.String())
}

type DesignatorItemDereference struct {
	DesignatorItem
}

func (*DesignatorItemDereference) Children() Nodes   { return Nodes{} }
func (*DesignatorItemDereference) isDesignatorItem() {}
func (*DesignatorItemDereference) String() string    { return "^" }

//   ```
//   Number
//   ```

type NumberFactor struct {
	Factor
	Value string
}

func NewNumber(v string) *NumberFactor { return &NumberFactor{Value: v} }
func (*NumberFactor) Children() Nodes  { return Nodes{} }
func (*NumberFactor) isFactor()        {}
func (m *NumberFactor) String() string { return m.Value }

//   ```
//   String
//   ```
type StringFactor struct {
	Factor
	Value string
}

func NewString(v string) *StringFactor { return &StringFactor{Value: v} }
func (*StringFactor) Children() Nodes  { return Nodes{} }
func (*StringFactor) isFactor()        {}
func (m *StringFactor) String() string { return m.Value }

// Ninl

type Nil struct {
	Factor
}

func NewNil() *Nil           { return &Nil{} }
func (*Nil) Children() Nodes { return Nodes{} }
func (*Nil) isFactor()       {}
func (*Nil) String() string  { return "Nil" }

// Parentheses
type Parentheses struct { // Round brackets
	Factor
	Expression *Expression
}

func (m *Parentheses) Children() Nodes { return Nodes{m.Expression} }
func (*Parentheses) isFactor()         {}
func (m *Parentheses) String() string {
	return fmt.Sprintf("(%s)", m.Expression.String())
}

//   ```
//   NOT Factor
//   ```

type Not struct {
	Factor
}

func (m *Not) Children() Nodes { return Nodes{m.Factor} }
func (*Not) isFactor()         {}
func (m *Not) String() string {
	return fmt.Sprintf("Not %s", m.Factor.String())
}

// - SetConstructor
//   ```
//   '[' [SetElement ','...] ']'
//   ```
type SetConstructor struct {
	Factor
	SetElements SetElements
}

func (SetConstructor) isFactor()          {}
func (m *SetConstructor) Children() Nodes { return Nodes{m.SetElements} }
func (m *SetConstructor) String() string {
	return fmt.Sprintf("[%s]", m.SetElements.String())
}

type SetElements []*SetElement // Must implement Node and fmt.Stringer
func (s SetElements) Children() Nodes {
	r := make(Nodes, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func (s SetElements) String() string {
	r := make(ext.Strings, len(s))
	for idx, i := range s {
		r[idx] = i.String()
	}
	return strings.Join(r, " ")
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
	return &SetElement{Expression: expr}
}

func (m *SetElement) Children() Nodes {
	r := Nodes{m.Expression}
	if m.SubRangeEnd != nil {
		r = append(r, m.SubRangeEnd)
	}
	return r
}
func (m *SetElement) String() string {
	if m.SubRangeEnd == nil {
		return m.Expression.String()
	} else {
		return fmt.Sprintf("%s..%s", m.Expression.String(), m.SubRangeEnd.String())
	}
}

//   ```
//   TypeId '(' Expression ')'
//   ```
type TypeCast struct {
	Factor
	TypeId     *TypeId
	Expression Expression
}

func (*TypeCast) isFactor()         {}
func (m *TypeCast) Children() Nodes { return Nodes{m.TypeId, &m.Expression} }
func (m *TypeCast) String() string {
	return fmt.Sprintf("%s(%s)", m.TypeId.String(), m.Expression.String())
}
