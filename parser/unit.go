package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) IsUnitIdentifier(t *token.Token) bool {
	s := t.Value()
	decl := p.context.Get(s)
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(*ast.UsesClauseItem)
	// log.Printf("UnitContext.IsUnitIdentifier(%s) decl.Node: %T %+v", s, decl.Node, decl.Node)
	return ok
}

func (p *Parser) ParseQualIds() (ast.QualIds, error) {
	res := ast.QualIds{}
	separator := token.Symbol(',')
	err := p.Until(token.Not(separator), separator, func() error {
		qualId, err := p.ParseQualId()
		if err != nil {
			return err
		}
		res = append(res, qualId)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, nil
	} else {
		return res, nil
	}
}

func (p *Parser) ParseQualId() (*ast.QualId, error) {
	if _, err := p.Current(token.Some(token.Identifier)); err != nil {
		return nil, err
	}
	name1 := p.CurrentToken()
	if p.IsUnitIdentifier(name1) {
		if _, err := p.Next(token.Symbol('.')); err != nil {
			return nil, err
		}
		name2, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		unitDecl := p.context.Get(name1.Value())
		if unitDecl == nil {
			return nil, p.TokenErrorf("undefined unit %s", name1)
		}
		usesClauseItem, ok := unitDecl.Node.(*ast.UsesClauseItem)
		if !ok {
			return nil, p.TokenErrorf("%s is not a unit but was %T (%+v)", name1, unitDecl.Node, unitDecl.Node)
		}
		unit := usesClauseItem.Unit
		if unit == nil {
			return nil, p.TokenErrorf("%s is used in uses clause but not found", name1)
		}
		decl := unit.DeclMap.Get(name2.Value())
		if decl == nil {
			return nil, p.TokenErrorf("undefined identifier %s in unit %s", name2, name1.Value())
		}
		p.NextToken()
		return &ast.QualId{
			NamespaceId: &ast.IdentRef{Ident: ast.NewIdent(name1), Ref: unitDecl},
			Ident:       &ast.IdentRef{Ident: ast.NewIdent(name2), Ref: decl},
		}, nil
	} else {
		p.NextToken()
		return ast.NewQualId(nil, p.NewIdentRef(name1)), nil
	}
}
