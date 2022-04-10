package ast

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

type Parameter struct {
	IdentList IdentList
	Type      *ParameterType
	ConstExpr *ConstExpr
}
