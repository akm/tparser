package parser

import (
	"io/ioutil"
	"os"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Program struct {
	*ast.Program
	Units ast.Units
}

func ParseProgram(path string) (*Program, error) {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	decoder := japanese.ShiftJIS.NewDecoder()
	str, err := ioutil.ReadAll(transform.NewReader(fp, decoder))
	if err != nil {
		return nil, err
	}

	runes := []rune(string(str))

	// absPath, err := filepath.Abs(path)
	// if err != nil {
	// 	return nil, err
	// }

	ctx := NewProgramContext(path)
	p := NewParser(&runes, ctx)
	p.NextToken()
	res, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	return &Program{
		Program: res,
		Units:   ctx.Units,
	}, nil
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("PROGRAM")); err != nil {
		return nil, err
	}
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Program{
		Path:  p.context.GetPath(),
		Ident: p.NewIdent(ident),
	}
	if _, err := p.Next(token.Symbol(';')); err != nil {
		return nil, err
	}
	p.NextToken()
	block, err := p.ParseProgramBlock()
	if err != nil {
		return nil, err
	}
	res.ProgramBlock = block
	if _, err := p.Current(token.Symbol('.')); err != nil {
		return nil, err
	}
	p.context.Set(res)
	return res, nil
}

func (p *Parser) ParseProgramBlock() (*ast.ProgramBlock, error) {
	res := &ast.ProgramBlock{}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("USES")) {
		uses, err := p.ParseUsesClause()
		if err != nil {
			return nil, err
		}
		res.UsesClause = uses
		p.NextToken()

		ctx, ok := p.context.(*ProgramContext)
		if !ok {
			panic(errors.Errorf("Something wrong. context is not ProjectContext"))
		}

		if err := p.LoadUnits(ctx, uses); err != nil {
			return nil, err
		}
	}
	block, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}
	res.Block = block
	return res, nil
}

func (p *Parser) LoadUnits(ctx *ProgramContext, uses ast.UsesClause) error {
	loaders := UnitLoaders{}
	for _, unitRef := range uses {
		path := unitRef.EffectivePath()
		if path != "" {
			loaders = append(loaders, NewUnitLoader(NewUnitContext(ctx, path)))
		}
	}

	for _, loader := range loaders {
		if err := loader.LoadFile(); err != nil {
			return err
		}
		if err := loader.ProcessIdentAndIntfUses(); err != nil {
			return err
		}
		ctx.AddUnit(loader.Unit)
	}

	sortedLoaders, err := loaders.Sort()
	if err != nil {
		return err
	}

	for _, loader := range sortedLoaders {
		if err := loader.ProcessIntfBody(); err != nil {
			return err
		}
	}

	declMaps := []astcore.DeclMap{ctx.DeclMap}
	declMaps = append(declMaps, loaders.DeclarationMaps()...)
	ctx.DeclMap = astcore.NewCompositeDeclarationMap(declMaps...)

	for _, loader := range sortedLoaders {
		if err := loader.ProcessImplAndInit(); err != nil {
			return err
		}
	}

	units := loaders.Units() // Don't use sortedLoaders for this
	for _, u := range units {
		ctx.DeclMap.Set(u)
	}

	return nil
}
