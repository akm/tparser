package ast

import (
	"strings"

	"github.com/akm/tparser/ext"
)

func IsStringTypeName(w string) bool {
	return stringTypeNames.Include(strings.ToUpper(w))
}

var stringTypeNames = ext.Strings{
	"STRING",
	"ANSISTRING",
	"WIDESTRING",
}.Set()

func (*StringType) isType() {}

type StringType struct {
	Name   string
	Length *ConstExpr
}
