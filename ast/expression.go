package ast

import "github.com/pkg/errors"

// - Expression
//   ```
//   SimpleExpression [RelOp SimpleExpression]...
//   ```
type Expression struct {
	SimpleExpression
	RelOpSimpleExpressions []*RelOpSimpleExpression
}

func NewExpression(arg interface{}) *Expression {
	switch v := arg.(type) {
	case *Expression:
		return v
	case *SimpleExpression:
		return &Expression{SimpleExpression: *v}
	default:
		return &Expression{SimpleExpression: *NewSimpleExpression(arg)}
	}
}

type ExprList []*Expression

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
	SimpleExpression
}

// - SimpleExpression
//   ```
//   ['+' | '-'] Term [AddOp Term]...
//   ```
type SimpleExpression struct {
	UnaryOp *string //  '+' | '-' or nil
	Term
	AddOpTerms []*AddOpTerm
}

func NewSimpleExpression(arg interface{}) *SimpleExpression {
	switch v := arg.(type) {
	case *SimpleExpression:
		return v
	case *Term:
		return &SimpleExpression{Term: *v}
	default:
		return &SimpleExpression{Term: *NewTerm(arg)}
	}
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
	Term
}

// - Term
//   ```
//   Factor [MulOp Factor]...
//   ```
type Term struct {
	Factor
	MulOpFactors []*MulOpFactor
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
}

// Designator ['(' ExprList ')']
func (*DesignatorFactor) isFactor() {}

type DesignatorFactor struct {
	Designator
	ExprList ExprList
}

func NewDesignatorFactor(arg interface{}) *DesignatorFactor {
	switch v := arg.(type) {
	case *DesignatorFactor:
		return v
	default:
		return &DesignatorFactor{Designator: *NewDesignator(arg)}
	}
}

// '@' Designator
func (*Address) isFactor() {}

type Address struct {
	Designator
}

// - Designator
//   ```
//   QualId ['.' Ident | '[' ExprList ']' | '^']...
//   ```
type Designator struct {
	QualId
	Items []DesignatorItem
}

func NewDesignator(arg interface{}) *Designator {
	switch v := arg.(type) {
	case *Designator:
		return v
	case *QualId:
		return &Designator{QualId: *v}
	case Ident:
		return &Designator{QualId: QualId{Ident: v}}
	case string:
		return &Designator{QualId: QualId{Ident: *NewIdent(v)}}
	default:
		panic(errors.Errorf("Unsupported type %T for NewDesignator", arg))
	}
}

type DesignatorItem interface {
	isDesignatorItem()
}

func (DesignatorItemIdent) isDesignatorItem() {}

type DesignatorItemIdent Ident

func NewDesignatorItemIdent(v interface{}) *DesignatorItemIdent {
	r := DesignatorItemIdent(*NewIdent(v))
	return &r
}

func (DesignatorItemExprList) isDesignatorItem() {}

type DesignatorItemExprList ExprList

func (*DesignatorItemDereference) isDesignatorItem() {}

type DesignatorItemDereference struct{}

func (ValueFactor) isFactor() {}

type ValueFactor struct {
	Value string
}

type Number struct{ ValueFactor }

func NewNumber(v string) *Number { return &Number{ValueFactor: ValueFactor{Value: v}} }

type String struct{ ValueFactor }

func NewString(v string) *String { return &String{ValueFactor: ValueFactor{Value: v}} }

// Ninl
func (*Nil) isFactor() {}

type Nil struct{}

func NewNil() *Nil { return &Nil{} }

// Parentheses
func (*Parentheses) isFactor() {}

type Parentheses struct { // Round brackets
	Expression Expression
}

func (*Not) isFactor() {}

type Not struct {
	Factor Factor
}

func (SetConstructor) isFactor() {}

// - SetConstructor
//   ```
//   '[' [SetElement ','...] ']'
//   ```
type SetConstructor struct {
	SetElements []*SetElement
}

// - SetElement
//   ```
//   Expression ['..' Expression]
//   ```
type SetElement struct {
	Expression
	SubRangeEnd *Expression
}

func NewSetElement(expr *Expression) *SetElement {
	return &SetElement{
		Expression: *expr,
	}
}

func (*TypeCast) isFactor() {}

type TypeCast struct {
	TypeId     *TypeId
	Expression Expression
}
