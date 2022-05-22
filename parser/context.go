package parser

import (
	"github.com/akm/tparser/parser/pcontext"
)

type (
	Context          = pcontext.Context
	ProgramContext   = pcontext.ProgramContext
	UnitContext      = pcontext.UnitContext
	StackableContext = pcontext.StackableContext
)

var (
	NewProgramContext   = pcontext.NewProgramContext
	NewUnitContext      = pcontext.NewUnitContext
	NewStackableContext = pcontext.NewStackableContext
)

func NewContext(args ...interface{}) Context {
	return NewProgramContext(args...)
}
