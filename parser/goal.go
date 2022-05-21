package parser

import (
	"github.com/akm/tparser/ast"
)

func (p *Parser) ParseGoal() (ast.Goal, error) {
	token := p.NextToken()
	switch token.Value() {
	case "PROGRAM":
		return p.ParseProgram()
	// case "PACKAGE":
	// 	return p.ParsePackage()
	// case "LIBRARY":
	// 	return p.ParseLibrary()
	case "UNIT":
		return p.ParseUnit()
	default:
		// 	return p.ParseProgram() // PROGRAM is optional word
		return nil, p.TokenErrorf("unexpected token: %s", token)
	}
}
