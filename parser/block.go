package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseBlock() (*ast.Block, error) {
	res := &ast.Block{}
	if declSections, err := p.ParseDeclSections(); err != nil {
		return nil, err
	} else if declSections != nil {
		res.DeclSections = declSections
	}

	if exportStmts, err := p.ParseExportsStmts(); err != nil {
		return nil, err
	} else if exportStmts != nil {
		res.ExportsStmts1 = exportStmts
	}

	switch p.CurrentToken().Value() {
	case "BEGIN":
		compoundStmt, err := p.ParseCompoundStmt(true)
		if err != nil {
			return nil, err
		}
		res.Body = compoundStmt
	case "ASM":
		asmStmt, err := p.ParseAssemblerStatement()
		if err != nil {
			return nil, err
		}
		res.Body = asmStmt
	}

	if exportStmts, err := p.ParseExportsStmts(); err != nil {
		return nil, err
	} else if exportStmts != nil {
		res.ExportsStmts1 = exportStmts
	}

	return res, nil
}

func (p *Parser) ParseDeclSections() (ast.DeclSections, error) {
	res := ast.DeclSections{}
	for {
		if sect, err := p.ParseDeclSection(); err != nil {
			return nil, err
		} else if sect != nil {
			res = append(res, sect)
		} else {
			break
		}
		if p.CurrentToken().Is(token.Symbol(';')) {
			p.NextToken()
		}
	}
	if len(res) == 0 {
		return nil, nil
	} else {
		return res, nil
	}
}

func (p *Parser) ParseDeclSection() (ast.DeclSection, error) {
	if sect, err := p.ParseLabelDeclSection(); err != nil {
		return nil, err
	} else if sect != nil {
		return sect, nil
	} else if sect, err := p.ParseConstSection(false); err != nil {
		return nil, err
	} else if sect != nil {
		return sect, nil
	} else if sect, err := p.ParseTypeSection(false); err != nil {
		return nil, err
	} else if sect != nil {
		return sect, nil
	} else if sect, err := p.ParseVarSection(false); err != nil {
		return nil, err
	} else if sect != nil {
		return sect, nil
	} else if sect, err := p.ParseThreadVarSection(false); err != nil {
		return nil, err
	} else if sect != nil {
		return sect, nil
	} else if sect, err := p.ParseProcedureDeclSection(); err != nil {
		return nil, err
	} else if sect != nil {
		return sect, nil
	}
	return nil, nil
}

func (p *Parser) ParseExportsStmts() (ast.ExportsStmts, error) {
	res := ast.ExportsStmts{}
	for {
		if stmt, err := p.ParseExportsStmt(false); err != nil {
			return nil, err
		} else if stmt != nil {
			res = append(res, stmt)
		} else {
			break
		}
	}
	if len(res) == 0 {
		return nil, nil
	} else {
		return res, nil
	}
}

func (p *Parser) ParseExportsStmt(required bool) (*ast.ExportsStmt, error) {
	kw := token.ReservedWord.HasKeyword("EXPORTS")
	if required {
		if _, err := p.Current(kw); err != nil {
			return nil, err
		}
	} else {
		if !p.CurrentToken().Is(kw) {
			return nil, nil
		}
	}
	p.NextToken()

	items := []*ast.ExportsItem{}

	if err := p.Until(token.Symbol(';'), token.Symbol(','), func() error {
		t, err := p.Current(token.Identifier)
		if err != nil {
			return err
		}
		item := &ast.ExportsItem{Ident: p.NewIdent(t)}
		t2 := p.NextToken()
		if t2.Is(token.Directives("NAME")) {
			p.NextToken()
			constExpr, err := p.ParseConstExpr()
			if err != nil {
				return err
			}
			item.Name = constExpr
		} else if t2.Is(token.Directives("INDEX")) {
			p.NextToken()
			constExpr, err := p.ParseConstExpr()
			if err != nil {
				return err
			}
			item.Index = constExpr
		}
		items = append(items, item)
		return nil
	}); err != nil {
		return nil, err
	}

	p.NextToken()
	return &ast.ExportsStmt{ExportsItems: items}, nil
}

func (p *Parser) ParseLabelDeclSection() (*ast.LabelDeclSection, error) {
	if !p.CurrentToken().Is(token.ReservedWord.HasKeyword("LABEL")) {
		return nil, nil
	}
	t, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	if _, err := p.Next(token.Symbol(';')); err != nil {
		return nil, err
	}
	p.NextToken()
	r := &ast.LabelDeclSection{LabelId: ast.NewLabelId(p.NewIdent(t))}
	p.context.Set(r)
	return r, nil
}
