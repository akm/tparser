package parser

import (
	"fmt"
	"strings"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseClassType() (ast.ClassType, error) {
	defer p.TraceMethod("Parser.ParseClassType")()

	if _, err := p.Current(token.ReservedWord.HasKeyword("CLASS")); err != nil {
		return nil, err
	}
	p.NextToken()
	if p.CurrentToken().Is(token.Symbol(';')) {
		return &ast.ForwardDeclaredClassType{}, nil
	}

	res := &ast.CustomClassType{}
	if heritage, err := p.ParseClassHeritage(); err != nil {
		return nil, err
	} else {
		res.Heritage = heritage
	}
	if memebrs, err := p.ParseClassMemberSections(res); err != nil {
		return nil, err
	} else {
		res.Members = memebrs
	}

	return nil, nil
}

func (p *Parser) ParseClassHeritage() (ast.ClassHeritage, error) {
	defer p.TraceMethod("Parser.ParseClassHeritage")()

	if !p.CurrentToken().Is(token.Symbol('(')) {
		return nil, nil
	}
	p.NextToken()

	res := ast.ClassHeritage{}
	if err := p.Until(token.Symbol(')'), token.Symbol(','), func() error {
		typeId, err := p.ParseTypeId()
		if err != nil {
			return err
		}
		res = append(res, typeId)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseClassMemberSections(classType *ast.CustomClassType) (ast.ClassMemberSections, error) {
	defer p.TraceMethod("Parser.ParseClassMemberSections")()

	res := ast.ClassMemberSections{}
	if err := p.Until(token.ReservedWord.HasKeyword("END"), nil, func() error {
		sect, err := p.ParseClassMemberSection(classType)
		if err != nil {
			return err
		}
		res = append(res, sect)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseClassMemberSection(classType *ast.CustomClassType) (*ast.ClassMemberSection, error) {
	defer p.TraceMethod("Parser.ParseClassMemberSection")()

	res := &ast.ClassMemberSection{}
	if t0, err := p.Current(token.Identifier); err != nil {
		return nil, err
	} else {
		switch strings.ToUpper(t0.Value()) {
		case "PRIVATE":
			res.Visibility = ast.CvPrivate
			p.NextToken()
		case "PROTECTED":
			res.Visibility = ast.CvProtected
			p.NextToken()
		case "PUBLIC":
			res.Visibility = ast.CvPblic
			p.NextToken()
		case "PUBLISHED":
			res.Visibility = ast.CvPublished
			p.NextToken()
		default:
			res.Visibility = ast.CvDefault
		}
	}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("END")) {
		return res, nil
	}

	fieldList, err := p.ParseClassFieldList()
	if err != nil {
		return nil, err
	}
	res.ClassFieldList = fieldList

	methodList, err := p.ParseClassMethodList()
	if err != nil {
		return nil, err
	}
	res.ClassMethodList = methodList

	propList, err := p.ParseClassPropertyList(classType)
	if err != nil {
		return nil, err
	}
	res.ClassPropertyList = propList

	return res, nil
}

var (
	propertyBreak = token.Some(
		token.ReservedWord.HasKeyword("END"),
	)
	methodBreak = token.Some(
		token.ReservedWord.HasKeyword("PROPERTY"),
		propertyBreak,
	)
	fieldListBreak = token.Some(
		token.ReservedWord.HasKeyword("FUNCTION"),
		token.ReservedWord.HasKeyword("PROCEDURE"),
		token.ReservedWord.HasKeyword("CONSTRUCTOR"),
		token.ReservedWord.HasKeyword("DESTRUCTOR"),
		methodBreak,
	)
)

func (p *Parser) ParseClassFieldList() (ast.ClassFieldList, error) {
	defer p.TraceMethod("Parser.ParseClassFieldList")()

	res := ast.ClassFieldList{}
	if err := p.Until(fieldListBreak, token.Symbol(';'), func() error {
		if fieldListBreak.Predicate(p.CurrentToken()) {
			return QuitUntil
		}
		field, err := p.ParseClassField()
		if err != nil {
			return err
		}
		res = append(res, field)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseClassField() (*ast.ClassField, error) {
	defer p.TraceMethod("Parser.ParseClassField")()

	identList, err := p.ParseIdentList(':')
	if err != nil {
		return nil, err
	}

	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res := &ast.ClassField{IdentList: *identList, Type: typ}
	return res, nil
}

func (p *Parser) ParseClassMethodList() (ast.ClassMethodList, error) {
	defer p.TraceMethod("Parser.ParseClassMethodList")()

	res := ast.ClassMethodList{}
	if err := p.Until(methodBreak, token.Symbol(';'), func() error {
		if methodBreak.Predicate(p.CurrentToken()) {
			return QuitUntil
		}
		method, err := p.ParseClassMethod()
		if err != nil {
			return err
		}
		res = append(res, method)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseClassMethod() (*ast.ClassMethod, error) {
	defer p.TraceMethod("Parser.ParseClassMethod")()

	res := &ast.ClassMethod{}
	t0, err := p.Current(token.ReservedWord)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(t0.Value()) {
	case "FUNCTION", "PROCEDURE":
		heading, err := p.ParseFunctionHeading()
		if err != nil {
			return nil, err
		}
		res.Heading = heading
	case "CONSTRUCTOR":
		heading, err := p.ParseConstructorHeading()
		if err != nil {
			return nil, err
		}
		res.Heading = heading
	case "DESTRUCTOR":
		heading, err := p.ParseDestructorHeading()
		if err != nil {
			return nil, err
		}
		res.Heading = heading
	default:
		return nil, fmt.Errorf("unexpected token for method: %s", t0)
	}

	directiveList, err := p.ClassMethodDirectiveList()
	if err != nil {
		return nil, err
	}
	res.Directives = directiveList

	return res, nil
}

func (p *Parser) ClassMethodDirectiveList() (ast.ClassMethodDirectiveList, error) {
	defer p.TraceMethod("Parser.ClassMethodDirectiveList")()

	res := ast.ClassMethodDirectiveList{}
	if err := p.Until(methodBreak, token.Symbol(';'), func() error {
		w := p.CurrentToken().Value()
		if ast.ClassMethodDirectives.Include(w) {
			res = append(res, ast.ClassMethodDirective(w))
			return nil
		} else {
			return QuitUntil
		}
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseConstructorHeading() (*ast.ConstructorHeading, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("CONSTRUCTOR")); err != nil {
		return nil, err
	}
	res := &ast.ConstructorHeading{}

	t0 := p.NextToken()
	res.Ident = ast.NewIdent(t0)

	p.NextToken()
	if p.CurrentToken().Is(token.Symbol('(')) {
		params, err := p.ParseFormalParameters('(', ')')
		if err != nil {
			return nil, err
		}
		res.FormalParameters = params
	}
	return res, nil
}

func (p *Parser) ParseDestructorHeading() (*ast.DestructorHeading, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("DESTRUCTOR")); err != nil {
		return nil, err
	}
	res := &ast.DestructorHeading{}

	t0 := p.NextToken()
	res.Ident = ast.NewIdent(t0)
	p.NextToken()
	return res, nil
}

func (p *Parser) ParseClassPropertyList(classType *ast.CustomClassType) (ast.ClassPropertyList, error) {
	res := ast.ClassPropertyList{}
	if err := p.Until(methodBreak, token.Symbol(';'), func() error {
		if methodBreak.Predicate(p.CurrentToken()) {
			return QuitUntil
		}
		prop, err := p.ParseClassProperty(classType)
		if err != nil {
			return err
		}
		res = append(res, prop)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseClassProperty(classType *ast.CustomClassType) (*ast.ClassProperty, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("PROPERTY")); err != nil {
		return nil, err
	}
	res := &ast.ClassProperty{}

	t0 := p.NextToken()
	res.Ident = ast.NewIdent(t0)
	p.NextToken()

	intf, err := p.ParsePropertyInterface()
	if err != nil {
		return nil, err
	}
	res.Interface = intf

	if strings.ToUpper(p.CurrentToken().Value()) == "INDEX" {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		res.Index = expr
	}

	//   [READ Ident]
	if strings.ToUpper(p.CurrentToken().Value()) == "READ" {
		t := p.NextToken()
		if decl := classType.FindMemberDecl(t.Value()); decl != nil {
			res.Read = ast.NewIdentRef(ast.NewIdent(t), decl)
		} else {
			return nil, fmt.Errorf("unknown member %s", t.Value())
		}
		p.NextToken()
	}

	//   [WRITE Ident]
	if strings.ToUpper(p.CurrentToken().Value()) == "WRITE" {
		t := p.NextToken()
		if decl := classType.FindMemberDecl(t.Value()); decl != nil {
			res.Write = ast.NewIdentRef(ast.NewIdent(t), decl)
		} else {
			return nil, fmt.Errorf("unknown member %s", t.Value())
		}
		p.NextToken()
	}

	//   [STORED (Ident | Constant)]
	if strings.ToUpper(p.CurrentToken().Value()) == "STORED" {
		t := p.NextToken()
		tVal := strings.ToLower(t.Value())
		if decl := classType.FindMemberDecl(tVal); decl != nil {
			res.Stored = &ast.PropertyStoredSpecifier{IdentRef: ast.NewIdentRef(ast.NewIdent(t), decl)}
		} else if tVal == "true" {
			v := true
			res.Stored = &ast.PropertyStoredSpecifier{Constant: &v}
		} else if tVal == "false" {
			v := false
			res.Stored = &ast.PropertyStoredSpecifier{Constant: &v}
		} else {
			return nil, fmt.Errorf("unknown member %s", t.Value())
		}
		p.NextToken()
	}

	//   [(DEFAULT ConstExpr) | NODEFAULT]
	if strings.ToUpper(p.CurrentToken().Value()) == "DEFAULT" {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		res.Default = &ast.PropertyDefaultSpecifier{Value: expr}
	} else if strings.ToUpper(p.CurrentToken().Value()) == "NODEFAULT" {
		v := true
		res.Default = &ast.PropertyDefaultSpecifier{NoDefault: &v}
	}

	//   [IMPLEMENTS TypeId]
	if strings.ToUpper(p.CurrentToken().Value()) == "IMPLEMENTS" {
		typeId, err := p.ParseTypeId()
		if err != nil {
			return nil, err
		}
		res.Implements = typeId
	}

	//   [PortabilityDirective]
	//    TODO

	return res, nil
}

func (p *Parser) ParsePropertyInterface() (*ast.PropertyInterface, error) {
	res := &ast.PropertyInterface{}
	if p.CurrentToken().Is(token.Symbol('[')) {
		params, err := p.ParseFormalParameters('[', ']')
		if err != nil {
			return nil, err
		}
		res.Parameters = params
	}
	if _, err := p.Current(token.Symbol(':')); err != nil {
		return nil, err
	}
	p.NextToken()

	typeId, err := p.ParseTypeId()
	if err != nil {
		return nil, err
	}
	res.Type = typeId
	return res, nil
}
