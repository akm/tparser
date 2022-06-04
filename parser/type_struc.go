package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
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
	case "SET":
		r, err := p.ParseSetType()
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
		indexes := []ast.OrdinalType{}
		if err := p.Until(token.Symbol(']'), token.Symbol(','), func() error {
			t0 := p.CurrentToken()
			typ, err := p.ParseType()
			if err != nil {
				return err
			}
			if ordinalType, ok := typ.(ast.OrdinalType); !ok {
				return errors.Errorf("Expected OrdinalType, got %T at %s", typ, p.PlaceString(t0))
			} else if !ordinalType.IsOrdinalType() {
				return errors.Errorf("Expected OrdinalType, got %T at %s", typ, p.PlaceString(t0))
			} else {
				indexes = append(indexes, ordinalType)
			}
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

func (p *Parser) ParseSetType() (*ast.SetType, error) {
	r := &ast.SetType{Packed: false}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("PACKED")) {
		r.Packed = true
		p.NextToken()
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("SET")); err != nil {
		return nil, p.TokenErrorf("Expected SET, got %s", p.CurrentToken())
	}
	p.NextToken()
	if _, err := p.Current(token.ReservedWord.HasKeyword("OF")); err != nil {
		return nil, p.TokenErrorf("Expected OF, got %s", p.CurrentToken())
	}
	p.NextToken()

	t0 := p.CurrentToken()
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	if ordinalType, ok := typ.(ast.OrdinalType); !ok {
		return nil, errors.Errorf("Expected OrdinalType, got %T at %s", typ, p.PlaceString(t0))
	} else if !ordinalType.IsOrdinalType() {
		return nil, errors.Errorf("Expected OrdinalType, got %T at %s", typ, p.PlaceString(t0))
	} else {
		r.OrdinalType = ordinalType
	}
	return r, nil
}
