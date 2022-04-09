package ast

type Directive string

const (
	DrCdecl       Directive = "CDECL"
	DrRegister    Directive = "REGISTER"
	DrDynamic     Directive = "DYNAMIC"
	DrVirtual     Directive = "VIRTUAL"
	DrExport      Directive = "EXPORT"
	DrExternal    Directive = "EXTERNAL" // `external` | `external stringConstant;` | `external stringConstant1 name stringConstant2;` | `external stringConstant index integerConstant;`
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
