package parser

import (
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Context struct {
	unitIdentifiers ext.Strings
}

func NewContext(args ...interface{}) *Context {
	var unitIdentifiers ext.Strings
	for _, arg := range args {
		switch v := arg.(type) {
		case ext.Strings:
			unitIdentifiers = v
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewContext", arg, arg))
		}
	}
	if unitIdentifiers == nil {
		unitIdentifiers = ext.Strings{}
	}
	return &Context{
		unitIdentifiers: unitIdentifiers,
	}
}

func (c *Context) IsUnitIdentifier(token *token.Token) bool {
	return c.unitIdentifiers.Include(token.Value())
}
