package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) IsUnitIdentifier() bool {
	return p.context.IsUnitIdentifier(p.CurrentToken())
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
	if p.context.IsUnitIdentifier(name1) {
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
		if !IsUsesClauseItem(unitDecl) {
			return nil, p.TokenErrorf("%s is not a unit", name1)
		}
		unit := unitDecl.Node.(*ast.Unit)
		decl := unit.DeclarationMap.Get(name2.Value())
		if decl == nil {
			return nil, p.TokenErrorf("undefined identifier %s in unit %s", name2, name1.Value())
		}
		p.NextToken()
		return &ast.QualId{
			UnitId: &ast.IdentRef{Ident: ast.NewIdent(name1), Ref: unitDecl},
			Ident:  &ast.IdentRef{Ident: ast.NewIdent(name2), Ref: decl},
		}, nil
	} else {
		p.NextToken()
		return ast.NewQualId(nil, p.NewIdentRef(name1)), nil
	}
}
