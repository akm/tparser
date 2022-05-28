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
	p := NewProgramParser(ctx)
	p.SetText(&runes)
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

type ProgramParser struct {
	*Parser
	context *ProgramContext
}

func NewProgramParser(ctx *ProgramContext) *ProgramParser {
	return &ProgramParser{Parser: NewParser(ctx), context: ctx}
}

func (p *ProgramParser) ParseProgram() (*ast.Program, error) {
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
	if err := p.context.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProgramParser) ParseProgramBlock() (*ast.ProgramBlock, error) {
	res := &ast.ProgramBlock{}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("USES")) {
		uses, err := p.ParseUsesClause()
		if err != nil {
			return nil, err
		}
		res.UsesClause = uses
		p.NextToken()

		if err := p.LoadUnits(p.context, uses); err != nil {
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

func (p *ProgramParser) LoadUnits(ctx *ProgramContext, uses ast.UsesClause) error {
	parsers := UnitParsers{}
	for _, unitRef := range uses {
		path := unitRef.EffectivePath()
		if path != "" {
			parsers = append(parsers, NewUnitParser(NewUnitContext(ctx, path)))
		}
	}

	for _, loader := range parsers {
		if err := loader.LoadFile(); err != nil {
			return err
		}
		if err := loader.ProcessIdentAndIntfUses(); err != nil {
			return err
		}
		ctx.AddUnit(loader.Unit)
	}

	sortedLoaders, err := parsers.Sort()
	if err != nil {
		return err
	}

	for _, loader := range sortedLoaders {
		if err := loader.ProcessIntfBody(); err != nil {
			return err
		}
	}

	declMaps := []astcore.DeclMap{ctx.DeclMap}
	declMaps = append(declMaps, parsers.DeclarationMaps()...)
	ctx.DeclMap = astcore.NewCompositeDeclarationMap(declMaps...)

	for _, loader := range sortedLoaders {
		if err := loader.ProcessImplAndInit(); err != nil {
			return err
		}
	}

	units := parsers.Units() // Don't use sortedLoaders for this
	for _, u := range units {
		usesItem := uses.Find(u.Ident.Name)
		if usesItem == nil {
			return errors.Errorf("UsesClauseItem not found for %s", u.Ident.Name)
		}
		usesItem.Unit = u
	}

	return nil
}
