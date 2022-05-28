package ast

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

// - ExportedHeading
//   ```
//   ProcedureHeading ';' [Directive]
//   ```
//   ```
//   FunctionHeading ';' [Directive]
//   ```
//
//
type ExportedHeading struct {
	*FunctionHeading
	Directives      []Directive
	ExternalOptions *ExternalOptions
	astcore.DeclNode
	InterfaceDecl
}

func (*ExportedHeading) canBeInterfaceDecl() {}
func (m *ExportedHeading) Children() Nodes   { return Nodes{m.FunctionHeading} }
func (m *ExportedHeading) GetDeclNodes() astcore.DeclNodes {
	return astcore.DeclNodes{m}
}

type FunctionType uint

const (
	FtProcedure FunctionType = iota
	FtFunction
)

func (m *ExportedHeading) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}

// - FunctionHeading
//   ```
//   FUNCTION Ident [FormalParameters] ':' (SimpleType | STRING)
//   ```
//   (Actually ReturnType is not only SimpleType or STRING.
//   TypeId also can be ReturnType.)
// - ProcedureHeading
//   ```
//   PROCEDURE Ident [FormalParameters]
//   ```

type FunctionHeading struct {
	Type FunctionType
	*Ident
	FormalParameters FormalParameters
	ReturnType       Type
}

func (*FunctionHeading) isExportedHeading() {}
func (s FunctionHeading) Children() Nodes {
	r := Nodes{s.Ident}
	if s.FormalParameters != nil {
		r = append(r, s.FormalParameters)
	}
	if s.ReturnType != nil {
		r = append(r, s.ReturnType)
	}
	return r
}

// - FormalParameters
//   ```
//   '(' [FormalParm ';'...] ')'
//   ```
type FormalParameters []*FormalParm

func (s FormalParameters) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - FormalParm
//   ```
//   [VAR | CONST | OUT] Parameter
//   ```
type FormalParmOption string

var (
	FpoVar   FormalParmOption = "VAR"
	FpoConst FormalParmOption = "CONST"
	FpoOut   FormalParmOption = "OUT"
)

type FormalParm struct {
	Opt *FormalParmOption
	*Parameter
	astcore.DeclNode
}

func (m *FormalParm) Children() Nodes {
	return Nodes{m.Parameter}
}

func NewFormalParm(name interface{}, args ...interface{}) *FormalParm {
	switch len(args) {
	case 0:
		switch v := name.(type) {
		case Parameter:
			return &FormalParm{Parameter: &v}
		case *Parameter:
			return &FormalParm{Parameter: v}
		default:
			panic(errors.Errorf("invalid argument for NewFormalParm: %v, %v", name, args))
		}
	case 1:
		return &FormalParm{Parameter: NewParameter(name, args[0])}
	case 2:
		var opt *FormalParmOption
		switch v := args[1].(type) {
		case FormalParmOption:
			opt = &v
		case *FormalParmOption:
			opt = v
		case string:
			switch v {
			case "VAR":
				opt = &FpoVar
			case "CONST":
				opt = &FpoConst
			case "OUT":
				opt = &FpoOut
			default:
				panic(errors.Errorf("invalid FormalParam option %q for NewFormalParm", v))
			}
		}
		return &FormalParm{Opt: opt, Parameter: NewParameter(name, args[0])}
	default:
		panic(errors.Errorf("too many arguments for NewFormalParm: %v, %v", name, args))
	}
}

func (m *FormalParm) ToDeclarations() astcore.Decls {
	return astcore.NewDeclarations(m.IdentList, m)
}

// - Parameter
//   ```
//   IdentList [':' ([ARRAY OF] SimpleType | STRING | FILE)]
//   ```
//   (Parameter type is not only SimpleType, STRING or FILE.
//   TypeId also can be also.)
//   ```
//   Ident ':' SimpleType '=' ConstExpr
//   ```
type ParameterType struct {
	Type
	IsArray bool
}

func (m *ParameterType) Children() Nodes {
	return Nodes{m.Type}
}

func NewParameterType(arg interface{}) *ParameterType {
	if arg == nil {
		return &ParameterType{}
	}
	switch v := arg.(type) {
	case *ParameterType:
		return v
	case ParameterType:
		return &v
	case Type:
		return &ParameterType{Type: v}
	default:
		return &ParameterType{Type: NewTypeId(arg)}
	}
}
func NewArrayParameterType(arg interface{}) *ParameterType {
	r := NewParameterType(arg)
	r.IsArray = true
	return r
}

type Parameter struct {
	IdentList
	Type      *ParameterType
	ConstExpr *ConstExpr
}

func (m *Parameter) Children() Nodes {
	r := Nodes{m.IdentList}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	if m.ConstExpr != nil {
		r = append(r, m.ConstExpr)
	}
	return r
}

func NewParameter(name interface{}, typArg interface{}, args ...interface{}) *Parameter {
	var typ *ParameterType
	if typArg != nil {
		typ = NewParameterType(typArg)
	}
	r := &Parameter{IdentList: NewIdentList(name), Type: typ}
	if len(args) > 1 {
		panic(errors.Errorf("too many arguments for NewParameter: %v, %v", name, args))
	} else if len(args) == 1 {
		r.SetConstExpr(args[0])
	}
	return r
}

func (m *Parameter) SetConstExpr(arg interface{}) {
	switch v := arg.(type) {
	case *ConstExpr:
		m.ConstExpr = v
	case ConstExpr:
		m.ConstExpr = &v
	default:
		m.ConstExpr = NewConstExpr(arg)
	}
}
