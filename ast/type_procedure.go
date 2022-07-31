package ast

// This is WRONG because ProcedureHeading and FunctionHeading have ident.
// - ProcedureType 🔖
//   ```
//   (ProcedureHeading | FunctionHeading) [OF OBJECT]
//   ```
//
// Actual:
// ```
// (FUNCTION | PROCEDURE) [FormalParameters] [':' (TypeId)] [of object]
// ```

type ProcedureType struct {
	FunctionType     FunctionType
	FormalParameters FormalParameters
	ReturnType       *TypeId
	OfObject         bool
}

var _ Type = (*ProcedureType)(nil)

func (*ProcedureType) isType() {}
func (m ProcedureType) Children() Nodes {
	r := Nodes{}
	if m.FormalParameters != nil {
		r = append(r, m.FormalParameters)
	}
	if m.ReturnType != nil {
		r = append(r, m.ReturnType)
	}
	return r
}
