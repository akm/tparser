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
// - ProcedureHeading
//   ```
//   PROCEDURE Ident [FormalParameters]
//   ```

func (*FunctionHeading) isExportedHeading() {}

type FunctionHeading struct {
	Type             FunctionType
	Ident            Ident
	FormalParameters FormalParameters
	ReturnType       *FunctionHeadingReturnType
}

type FunctionHeadingReturnType struct {
	SimpleType SimpleType // allow nil
	TypeName   *string    // STRING or nil
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

const (
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
//   ```
//   Ident ':' SimpleType '=' ConstExpr
//   ```
type ParameterType struct {
	IsArray    bool
	SimpleType SimpleType // allow nil
	TypeName   *string    // STRING, FILE or nil
}

type Parameter struct {
	IdentList IdentList
	Type      *ParameterType
	ConstExpr *ConstExpr
}
