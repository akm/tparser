package parser

import (
	"github.com/akm/opparser/ast"
	"github.com/pkg/errors"
)

func (p *Parser) ParseFile() (ast.Goal, error) {
	token := p.next()
	switch token.Text() {
	// case "PROGRAM":
	// 	return p.ParseProgram()
	// case "PACKAGE":
	// 	return p.ParsePackage()
	// case "LIBRARY":
	// 	return p.ParseLibrary()
	case "UNIT":
		return p.ParseUnit()
	default:
		// 	return p.ParseProgram() // PROGRAM is optional word
		return nil, errors.Errorf("unexpected token: %v", token)
	}
}
