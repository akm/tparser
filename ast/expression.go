package ast

import "github.com/pkg/errors"

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

type RelOpSimpleExpression struct {
	RelOp string // '>' | '<' | '<=' | '>=' | '=' | '<>' | "IN" | "IS" | "AS"
	SimpleExpression
}

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

type AddOpTerm struct {
	AddOp string // '+' | '-' | "OR" | "XOR"
	Term
}

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

type MulOpFactor struct {
	MulOp string // '*' | '/' | "DIV" | "MOD", "AND", "SHL", "SHR"
	Factor
}

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

// QualId ['.' Ident | '[' ExprList ']' | '^']...
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
		return &Designator{QualId: QualId{Ident: Ident(v)}}
	default:
		panic(errors.Errorf("Unsupported type %T for NewDesignator", arg))
	}
}

type DesignatorItem interface {
	isDesignatorItem()
}

func (DesignatorItemIdent) isDesignatorItem() {}

type DesignatorItemIdent Ident

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

func (*SetConstructor) isFactor() {}

type SetConstructor struct {
	SetElements []*SetElement
}

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
