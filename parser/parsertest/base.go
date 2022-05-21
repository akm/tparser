package parsertest

import (
	"github.com/akm/tparser/parser"
)

func NewTestParser(text *[]rune, origArgs ...interface{}) *parser.Parser {
	args := []interface{}{}
	var ctx parser.Context
	for _, origArg := range origArgs {
		switch v := origArg.(type) {
		case parser.Context:
			ctx = v
		default:
			args = append(args, origArg)
		}
	}
	if ctx == nil {
		ctx = NewTestProgramContext()
	}
	return parser.NewParser(text, ctx, args...)
}

func NewTestProgramContext(args ...interface{}) *parser.ProgramContext {
	return parser.NewProgramContext(args...)
}
