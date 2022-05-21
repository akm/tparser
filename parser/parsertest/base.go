package parsertest

import (
	"github.com/akm/tparser/parser"
	"github.com/pkg/errors"
)

func NewTestParser(text *[]rune, origArgs ...interface{}) *parser.Parser {
	var ctx parser.Context
	for _, origArg := range origArgs {
		switch v := origArg.(type) {
		case parser.Context:
			ctx = v
		default:
			panic(errors.Errorf("Unsupported type %T (%v) is given for NewTestParser", origArg, origArg))
		}
	}
	if ctx == nil {
		ctx = NewTestUnitContext()
	}
	p := parser.NewParser(ctx)
	p.SetText(text)
	return p
}

func NewTestProgramParser(text *[]rune, origArgs ...interface{}) *parser.ProgramParser {
	var ctx *parser.ProgramContext
	for _, origArg := range origArgs {
		switch v := origArg.(type) {
		case *parser.ProgramContext:
			ctx = v
		default:
			panic(errors.Errorf("Unsupported type %T (%v) is given for NewTestParser", origArg, origArg))
		}
	}
	if ctx == nil {
		ctx = NewTestProgramContext()
	}
	return parser.NewProgramParser(text, ctx)
}

func NewTestUnitParser(text *[]rune, origArgs ...interface{}) *parser.UnitParser {
	var ctx *parser.UnitContext
	for _, origArg := range origArgs {
		switch v := origArg.(type) {
		case *parser.UnitContext:
			ctx = v
		default:
			panic(errors.Errorf("Unsupported type %T (%v) is given for NewTestParser", origArg, origArg))
		}
	}
	if ctx == nil {
		ctx = NewTestUnitContext()
	}
	return parser.NewUnitParser(text, ctx)
}

func NewTestProgramContext(args ...interface{}) *parser.ProgramContext {
	return parser.NewProgramContext(args...)
}

func NewTestUnitContext(args ...interface{}) *parser.UnitContext {
	return parser.NewUnitContext(NewTestProgramContext(), args...)
}
