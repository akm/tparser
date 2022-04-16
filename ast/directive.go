package ast

import "github.com/pkg/errors"

type Directive string

const (
	DrCdecl       Directive = "CDECL"
	DrRegister    Directive = "REGISTER"
	DrDynamic     Directive = "DYNAMIC"
	DrVirtual     Directive = "VIRTUAL"
	DrExport      Directive = "EXPORT"
	DrExternal    Directive = "EXTERNAL"
	DrNear        Directive = "NEAR"
	DrFar         Directive = "FAR"
	DrForward     Directive = "FORWARD"
	DrMessage     Directive = "MESSAGE" // MESSAGE ConstExpr
	DrOverride    Directive = "OVERRIDE"
	DrOverload    Directive = "OVERLOAD"
	DrPascal      Directive = "PASCAL"
	DrReintroduce Directive = "REINTRODUCE"
	DrSafecall    Directive = "SAFECALL"
	DrStdcall     Directive = "STDCALL"
	DrVarArgs     Directive = "VARARGS"
	DrLocal       Directive = "LOCAL"
	DrAbstract    Directive = "ABSTRACT"
)

// Directive external can be used like these
//   external;
//   external stringConstant;
//   external stringConstant1 name stringConstant2;
//   external stringConstant index integerConstant;
type ExternalOptions struct {
	LibraryName string
	Name        *string
	Index       *int
}

func NewExternalOptions(libraryName string, args ...interface{}) *ExternalOptions {
	if len(args) > 1 {
		panic(errors.Errorf("too many arguments for NewExternalOptions: %v, %v", libraryName, args))
	}
	r := &ExternalOptions{LibraryName: libraryName}
	if len(args) == 1 {
		switch v := args[0].(type) {
		case string:
			r.Name = &v
		case *string:
			r.Name = v
		case int:
			r.Index = &v
		case *int:
			r.Index = v
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewExternalOptions", args[0], args[0]))
		}
	}
	return r
}
