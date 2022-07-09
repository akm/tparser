package parser

import (
	"strconv"
	"strings"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

func (p *Parser) ParseExportedHeading() (*ast.ExportedHeading, error) {
	defer p.context.StackDeclMap()()

	var functionHeading *ast.FunctionHeading
	switch p.CurrentToken().Value() {
	case "PROCEDURE", "FUNCTION":
		var err error
		functionHeading, err = p.ParseFunctionHeading()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.Current(token.Symbol(';')); err != nil {
		return nil, err
	}
	r := &ast.ExportedHeading{FunctionHeading: functionHeading}
	if err := p.context.Set(r); err != nil {
		return nil, err
	}

	p.NextToken()
	if p.CurrentToken().Is(token.Directive) {
		directives, opts, err := p.ParseFunctionDirectives()
		if err != nil {
			return nil, err
		}
		r.Directives = directives
		r.ExternalOptions = opts
	}
	return r, nil
}

func (p *Parser) ParseFunctionDirectives() ([]ast.Directive, *ast.ExternalOptions, error) {
	directives := []ast.Directive{}
	var opts *ast.ExternalOptions
	for {
		t := p.CurrentToken()
		if !t.Is(token.Directive) {
			break
		}
		dir := ast.Directive(strings.ToUpper(t.Value()))
		directives = append(directives, dir)
		switch dir {
		case ast.DrExternal:
			extOpts, err := p.ParseExternalOptions()
			if err != nil {
				return nil, nil, err
			}
			opts = extOpts
		default:
			p.NextToken()
		}
		if _, err := p.Current(token.Symbol(';')); err != nil {
			return nil, nil, err
		}
		p.NextToken()
	}
	return directives, opts, nil
}

func (p *Parser) ParseExternalOptions() (*ast.ExternalOptions, error) {
	if _, err := p.Current(token.Directives("EXTERNAL")); err != nil {
		return nil, err
	}
	t := p.NextToken()
	if t.Is(token.Symbol(';')) {
		return nil, nil
	}
	r := &ast.ExternalOptions{}
	f1, err := p.ParseStringFactor(t, false)
	if err != nil {
		return nil, err
	}

	r.LibraryName = f1.Value
	if p.CurrentToken().Is(token.Symbol(';')) {
		return r, nil
	}
	{
		t := p.CurrentToken()
		if t.Is(token.Directives("NAME")) {
			t2 := p.NextToken()
			f2, err := p.ParseStringFactor(t2, false)
			if err != nil {
				return nil, err
			}
			r.Name = &f2.Value
		} else if t.Is(token.Directives("INDEX")) {
			t2 := p.NextToken()
			f2, err := p.ParseNumberFactor(t2, false)
			if err != nil {
				return nil, err
			}
			val, err := strconv.ParseInt(f2.Value, 10, 64)
			if err != nil {
				return nil, errors.Errorf("Invalid index constant for extenal index %s", f2.Value)
			}
			v := int(val)
			r.Index = &v
		} else {
			return nil, p.TokenErrorf("expects NAME or INDEX, but got %s (%s)", t, string(t.Raw()))
		}
	}

	return r, nil
}

func (p *Parser) ParseFunctionHeading() (*ast.FunctionHeading, error) {
	res := &ast.FunctionHeading{}

	t0 := p.CurrentToken()
	switch t0.Value() {
	case "FUNCTION":
		res.Type = ast.FtFunction
	case "PROCEDURE":
		res.Type = ast.FtProcedure
	default:
		return nil, p.TokenErrorf("expects FUNCTION or PROCEDURE, but got %s (%s)", t0, string(t0.Raw()))
	}

	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res.Ident = p.NewIdent(ident)
	t := p.NextToken()
	if t.Is(token.Symbol('(')) {
		formalParameters, err := p.ParseFormalParameters('(', ')')
		if err != nil {
			return nil, err
		}
		res.FormalParameters = formalParameters
	}
	if res.Type == ast.FtFunction {
		if _, err := p.Current(token.Symbol(':')); err != nil {
			return nil, err
		}
		p.NextToken()
		typ, err := p.ParseTypeId()
		if err != nil {
			return nil, err
		}
		res.ReturnType = typ
	}
	return res, nil
}

func (p *Parser) ParseFormalParameters(startRune, endRune rune) (ast.FormalParameters, error) {
	if _, err := p.Current(token.Symbol(startRune)); err != nil {
		return nil, err
	}
	p.NextToken()
	r := ast.FormalParameters{}
	if err := p.Until(token.Symbol(endRune), token.Symbol(';'), func() error {
		formalParm, err := p.ParseFormalParm()
		if err != nil {
			return err
		}
		r = append(r, formalParm)
		return nil
	}); err != nil {
		return nil, err
	}
	// // Until で以下はチェック済み
	// if _, err := p.Current(token.Symbol(')')); err != nil {
	// 	return nil, err
	// }
	p.NextToken()
	return r, nil
}

func (p *Parser) ParseFormalParm() (*ast.FormalParm, error) {
	r := &ast.FormalParm{}
	t := p.CurrentToken()
	if t.Is(token.ReservedWord) {
		switch t.Value() {
		case "VAR":
			r.Opt = &ast.FpoVar
		case "CONST":
			r.Opt = &ast.FpoConst
		case "OUT":
			r.Opt = &ast.FpoOut
		default:
			return nil, p.TokenErrorf("unexpected token %s", t)
		}
		p.NextToken()
	}
	parameter, err := p.ParseParameter()
	if err != nil {
		return nil, err
	}
	r.Parameter = parameter
	if err := p.context.Set(r); err != nil {
		return nil, err
	}
	return r, nil
}

var parameterTerminators = token.Some(
	token.Symbol(';'),
	token.Symbol(')'),
)
var parameterIdentListTerminators = token.Some(
	token.Symbol(':'),
	token.Symbol(';'),
	token.Symbol(')'),
)

func (p *Parser) ParseParameter() (*ast.Parameter, error) {
	identList, err := p.ParseIdentListBy(parameterIdentListTerminators)
	if err != nil {
		return nil, err
	}
	r := &ast.Parameter{
		IdentList: *identList,
	}
	if p.CurrentToken().Is(token.Symbol(':')) {
		parameterType := &ast.ParameterType{}
		p.NextToken()
		if p.CurrentToken().Is(token.ReservedWord.HasKeyword("ARRAY")) {
			if _, err := p.Next(token.ReservedWord.HasKeyword("OF")); err != nil {
				return nil, err
			}
			parameterType.IsArray = true
			p.NextToken()
		}
		if parameterType.IsArray && p.CurrentToken().Is(token.ReservedWord.HasKeyword("CONST")) {
			parameterType.Type = nil // array of const
			p.NextToken()
		} else {
			typ, err := p.ParseType()
			if err != nil {
				return nil, err
			}
			parameterType.Type = typ
		}
		r.Type = parameterType
	}
	if p.CurrentToken().Is(token.Symbol('=')) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		r.ConstExpr = expr
	}
	if _, err := p.Current(parameterTerminators); err != nil {
		return nil, err
	}
	return r, nil
}
