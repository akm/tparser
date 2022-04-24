package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewFormalParm(name interface{}, args ...interface{}) *ast.FormalParm {
	if sname, ok := name.(string); ok {
		name = NewIdent(sname)
	} else if nameslice, ok := name.([]string); ok {
		v := make([]interface{}, len(nameslice))
		for idx, i := range nameslice {
			v[idx] = NewIdent(i)
		}
		name = v
	}
	if len(args) > 0 {
		if arg0, ok := args[0].(string); ok {
			args[0] = NewTypeId(arg0)
		}
	}
	return ast.NewFormalParm(name, args...)
}

func NewParameterType(arg interface{}) *ast.ParameterType {
	if arg == nil {
		return &ast.ParameterType{}
	}
	switch v := arg.(type) {
	case string:
		return ast.NewParameterType(NewIdent(v))
	default:
		return ast.NewParameterType(arg)
	}
}

func NewArrayParameterType(arg interface{}) *ast.ParameterType {
	switch v := arg.(type) {
	case string:
		return ast.NewArrayParameterType(NewIdent(v))
	default:
		return ast.NewArrayParameterType(arg)
	}
}

func NewParameter(name interface{}, typArg interface{}, args ...interface{}) *ast.Parameter {
	if sname, ok := name.(string); ok {
		name = NewIdent(sname)
	} else if nameslice, ok := name.([]string); ok {
		v := make([]interface{}, len(nameslice))
		for idx, i := range nameslice {
			v[idx] = NewIdent(i)
		}
		name = v
	}
	if styp, ok := typArg.(string); ok {
		typArg = NewTypeId(styp)
	}
	if len(args) == 1 {
		arg := args[0]
		if sarg, ok := arg.(string); ok {
			args[0] = NewIdent(sarg)
		}
	}
	return ast.NewParameter(name, typArg, args...)
}
