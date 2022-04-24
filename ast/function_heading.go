package ast

import "github.com/pkg/errors"

// - ExportedHeading
//   ```
//   ProcedureHeading ';' [Directive]
//   ```
//   ```
//   FunctionHeading ';' [Directive]
//   ```
//
//
func (*ExportedHeading) canBeInterfaceDecl() {}

type ExportedHeading struct {
	FunctionHeading FunctionHeading
	Directives      []Directive
	ExternalOptions *ExternalOptions
}

type FunctionType uint

const (
	FtProcedure FunctionType = iota
	FtFunction
)

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

func (*FunctionHeading) isExportedHeading() {}

type FunctionHeading struct {
	Type             FunctionType
	Ident            Ident
	FormalParameters FormalParameters
	ReturnType       Type
}

// - FormalParameters
//   ```
//   '(' [FormalParm ';'...] ')'
//   ```
type FormalParameters []*FormalParm

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
	Parameter
}

func NewFormalParm(name interface{}, args ...interface{}) *FormalParm {
	switch len(args) {
	case 0:
		switch v := name.(type) {
		case Parameter:
			return &FormalParm{Parameter: v}
		case *Parameter:
			return &FormalParm{Parameter: *v}
		default:
			panic(errors.Errorf("invalid argument for NewFormalParm: %v, %v", name, args))
		}
	case 1:
		return &FormalParm{Parameter: *NewParameter(name, args[0])}
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
		return &FormalParm{Opt: opt, Parameter: *NewParameter(name, args[0])}
	default:
		panic(errors.Errorf("too many arguments for NewFormalParm: %v, %v", name, args))
	}
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
	IdentList IdentList
	Type      *ParameterType
	ConstExpr *ConstExpr
}

func NewParameter(name interface{}, typArg interface{}, args ...interface{}) *Parameter {
	var typ *ParameterType
	if typArg != nil {
		switch v := typArg.(type) {
		case ParameterType:
			typ = &v
		case *ParameterType:
			typ = v
		case Type:
			typ = &ParameterType{Type: v}
		default:
			typ = NewParameterType(typArg)
		}
	}
	r := &Parameter{IdentList: NewIdentList(name), Type: typ}
	if len(args) > 1 {
		panic(errors.Errorf("too many arguments for NewParameter: %v, %v", name, args))
	}
	if len(args) == 1 {
		switch v := args[0].(type) {
		case *ConstExpr:
			r.ConstExpr = v
		case ConstExpr:
			r.ConstExpr = &v
		default:
			r.ConstExpr = NewConstExpr(args[0])
		}
	}
	return r
}
