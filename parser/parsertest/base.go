package parsertest

import "github.com/akm/tparser/parser"

func NewTestParser(text *[]rune, args ...interface{}) *parser.Parser {
	return parser.NewParser(text, args...)
}
