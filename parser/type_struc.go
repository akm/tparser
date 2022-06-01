package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseStrucType() (ast.StrucType, error) {
	packed := false
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("PACKED")) {
		packed = true
		p.NextToken()
	}
	switch p.CurrentToken().Value() {
	case "ARRAY":
		r, err := p.ParseArrayType()
		if err != nil {
			return nil, err
		}
		if !r.Packed && packed {
			r.Packed = packed
		}
		return r, nil
	default:
		return nil, p.TokenErrorf("Unsupported StrucType token %s", p.CurrentToken())
	}
}

func (p *Parser) ParseArrayType() (*ast.ArrayType, error) {
	r := &ast.ArrayType{Packed: false}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("PACKED")) {
		r.Packed = true
		p.NextToken()
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("ARRAY")); err != nil {
		return nil, p.TokenErrorf("Expected ARRAY, got %s", p.CurrentToken())
	}
	p.NextToken()

	if p.CurrentToken().Is(token.Symbol('[')) {
		p.NextToken()
		indexes := []ast.Type{}
		if err := p.Until(token.Symbol(']'), token.Symbol(','), func() error {
			typ, err := p.ParseType()
			if err != nil {
				return err
			}
			indexes = append(indexes, typ)
			return nil
		}); err != nil {
			return nil, err
		}
		r.IndexTypes = indexes
		p.NextToken()
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("OF")); err != nil {
		return nil, p.TokenErrorf("Expected OF, got %s", p.CurrentToken())
	}
	p.NextToken()

	baseType, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	r.BaseType = baseType

	if p.NextToken().Is(token.ReservedWord.HasKeyword("PACKED")) {
		r.Packed = true
		p.NextToken()
	}

	return r, nil
}
