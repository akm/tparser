package ast

import (
	"strings"

	"github.com/akm/tparser/ext"
	"github.com/pkg/errors"
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

func NewStringType(name string, args ...interface{}) *StringType {
	if len(args) > 1 {
		panic(errors.Errorf("too many arguments for NewStringType: %v, %v", name, args))
	}
	var length *ConstExpr
	if len(args) == 1 {
		switch v := args[0].(type) {
		case ConstExpr:
			length = &v
		case *ConstExpr:
			length = v
		}
	}
	return &StringType{
		Name:   name,
		Length: length,
	}
}
