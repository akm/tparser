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

type ExportedHeading struct {
	FunctionHeading FunctionHeading
	Directive       []Directive
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

func NewFormalParm(name interface{}, typ interface{}, args ...interface{}) *FormalParm {
	var opt *FormalParmOption
	if len(args) > 1 {
		panic(errors.Errorf("too many arguments for NewFormalParm: %v, %v", name, args))
	} else if len(args) == 1 {
		switch v := args[0].(type) {
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
	}
	return &FormalParm{
		Opt:       opt,
		Parameter: *NewParameter(name, typ),
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

func NewParameter(name interface{}, typArg interface{}) *Parameter {
	var typ *ParameterType
	if typArg != nil {
		switch v := typArg.(type) {
		case ParameterType:
			typ = &v
		case *ParameterType:
			typ = v
		default:
			typ = NewParameterType(typArg)
		}
	}
	return &Parameter{IdentList: NewIdentList(name), Type: typ}
}
