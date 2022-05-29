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
	// log.Printf("Parser.IsUnitIdentifier(%s) decl.Node: %T %+v", s, decl.Node, decl.Node)
	return ok
}

func (p *Parser) IsNamespaceIdentifier(t *token.Token) bool {
	s := t.Value()
	decl := p.context.Get(s)
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(ast.Namespace)
	// log.Printf("Parser.IsNamespaceIdentifier(%s) decl.Node: %T %+v", s, decl.Node, decl.Node)
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
	if p.IsNamespaceIdentifier(name1) || p.IsUnitIdentifier(name1) {
		if _, err := p.Next(token.Symbol('.')); err != nil {
			return nil, err
		}
		name2, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		namespaceDecl := p.context.Get(name1.Value())
		if namespaceDecl == nil {
			return nil, p.TokenErrorf("undefined unit %s", name1)
		}
		var namespace ast.Namespace
		if usesClauseItem, ok := namespaceDecl.Node.(*ast.UsesClauseItem); ok {
			unit := usesClauseItem.Unit
			if unit == nil {
				return nil, p.TokenErrorf("%s is used in uses clause but not found", name1)
			}
			namespace = unit
		} else if program, ok := namespaceDecl.Node.(ast.Namespace); ok {
			namespace = program
		} else {
			return nil, p.TokenErrorf("%s is neither unit nor program but was %T", name1, namespaceDecl.Node)
		}

		decl := namespace.GetDeclMap().Get(name2.Value())
		if decl == nil {
			return nil, p.TokenErrorf("undefined identifier %s in unit %s", name2, name1.Value())
		}
		p.NextToken()
		return &ast.QualId{
			NamespaceId: &ast.IdentRef{Ident: ast.NewIdent(name1), Ref: namespaceDecl},
			Ident:       &ast.IdentRef{Ident: ast.NewIdent(name2), Ref: decl},
		}, nil
	} else {
		p.NextToken()
		return ast.NewQualId(nil, p.NewIdentRef(name1)), nil
	}
}
