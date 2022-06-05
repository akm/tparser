package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/log"
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
	case "SET":
		r, err := p.ParseSetType()
		if err != nil {
			return nil, err
		}
		if !r.Packed && packed {
			r.Packed = packed
		}
		return r, nil
	case "RECORD":
		r, err := p.ParseRecType()
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
			ordinalType, err := p.ParseTypeAsOrdinalType()
			if err != nil {
				return err
			}
			indexes = append(indexes, ordinalType)
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

	ordinalType, err := p.ParseTypeAsOrdinalType()
	if err != nil {
		return nil, err
	}
	r.OrdinalType = ordinalType
	return r, nil
}

func (p *Parser) ParseRecType() (*ast.RecType, error) {
	r := &ast.RecType{Packed: false}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("PACKED")) {
		r.Packed = true
		p.NextToken()
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("RECORD")); err != nil {
		return nil, p.TokenErrorf("Expected RECORD, got %s", p.CurrentToken())
	}
	p.NextToken()
	fieldList, err := p.ParseFieldList(token.ReservedWord.HasKeyword("END"))
	if err != nil {
		return nil, err
	}
	r.FieldList = fieldList

	if _, err := p.Current(token.ReservedWord.HasKeyword("END")); err != nil {
		return nil, err
	}

	return r, nil
}

func (p *Parser) ParseFieldList(terminator token.Predicator) (*ast.FieldList, error) {
	r := &ast.FieldList{}
	casePred := token.ReservedWord.HasKeyword("CASE")
	fieldDecls := ast.FieldDecls{}
	if err := p.Until(terminator, token.Symbol(';'), func() error {
		if p.CurrentToken().Is(casePred) || p.CurrentToken().Is(terminator) {
			return QuitUntil
		}
		fieldDecl, err := p.ParseFieldDecl(terminator)
		if err != nil {
			return err
		}
		log.Printf("fieldDecl: %s", fieldDecl.String())
		fieldDecls = append(fieldDecls, fieldDecl)
		return nil
	}); err != nil {
		return nil, err
	}
	r.FieldDecls = fieldDecls

	if p.CurrentToken().Is(casePred) {
		variantSection, err := p.ParseVariantSection()
		if err != nil {
			return nil, err
		}
		r.VariantSection = variantSection
	} else if !p.CurrentToken().Is(terminator) {
		p.NextToken()
	}

	return r, nil
}

func (p *Parser) ParseFieldDecl(terminator token.Predicator) (*ast.FieldDecl, error) {
	identList, err := p.ParseIdentList(':')
	if err != nil {
		return nil, err
	}
	r := &ast.FieldDecl{IdentList: *identList}
	p.NextToken()
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	r.Type = typ
	if p.CurrentToken().Is(terminator) {
		return r, nil
	}
	if _, err := p.Current(token.Symbol(';')); err != nil {
		return nil, err
	}
	return r, nil
}

func (p *Parser) ParseVariantSection() (*ast.VariantSection, error) {
	log.Printf("ParseVariantSection")
	if _, err := p.Current(token.ReservedWord.HasKeyword("CASE")); err != nil {
		return nil, p.TokenErrorf("Expected CASE, got %s", p.CurrentToken())
	}
	p.NextToken()
	r := &ast.VariantSection{}

	rollback := p.RollbackPoint()
	identToken, err := p.Current(token.Identifier)
	if err != nil {
		return nil, err
	}
	t := p.NextToken()
	if t.Is(token.Symbol(':')) {
		p.NextToken()
		r.Ident = ast.NewIdent(identToken)
	} else {
		rollback()
	}

	ordinalType, err := p.ParseTypeAsOrdinalType()
	if err != nil {
		return nil, err
	}
	r.TypeId = ordinalType

	if _, err := p.Current(token.ReservedWord.HasKeyword("OF")); err != nil {
		return nil, p.TokenErrorf("Expected OF, got %s", p.CurrentToken())
	}
	p.NextToken()

	recVariants := ast.RecVariants{}
	endPred := token.ReservedWord.HasKeyword("END")
	if err := p.Until(endPred, token.Symbol(';'), func() error {
		if p.CurrentToken().Is(endPred) {
			return QuitUntil
		}
		recVariant, err := p.ParseRecVariant()
		if err != nil {
			return err
		}
		recVariants = append(recVariants, recVariant)
		return nil
	}); err != nil {
		return nil, err
	}
	r.RecVariants = recVariants
	// p.NextToken() // Don't go to next token because ParseRecType check whether current token is END

	return r, nil
}

func (p *Parser) ParseRecVariant() (*ast.RecVariant, error) {
	constExprs := ast.ConstExprs{}
	if err := p.Until(token.Symbol(':'), token.Symbol(','), func() error {
		constExpr, err := p.ParseConstExpr()
		if err != nil {
			return err
		}
		constExprs = append(constExprs, constExpr)
		return nil
	}); err != nil {
		return nil, err
	}
	if _, err := p.Next(token.Symbol('(')); err != nil {
		return nil, err
	}
	p.NextToken()

	fieldList, err := p.ParseFieldList(token.Symbol(')'))
	if err != nil {
		return nil, err
	}
	p.NextToken() // Go to next token because ParseFieldList doesn't call NextToken which quits by terminator
	return &ast.RecVariant{
		ConstExprs: constExprs,
		FieldList:  fieldList,
	}, nil
}
