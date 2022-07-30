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
	if !p.CurrentToken().Is(token.Symbol(';')) {
		if _, err := p.ParseClassMemberSections(res); err != nil {
			return nil, err
		}
	}

	return res, nil
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
	p.NextToken()
	return res, nil
}

func (p *Parser) ParseClassMemberSections(classType *ast.CustomClassType) (ast.ClassMemberSections, error) {
	defer p.TraceMethod("Parser.ParseClassMemberSections")()

	classType.Members = ast.ClassMemberSections{}
	if err := p.Until(token.ReservedWord.HasKeyword("END"), nil, func() error {
		if _, err := p.ParseClassMemberSection(classType); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return classType.Members, nil
}

func (p *Parser) ParseClassMemberSection(classType *ast.CustomClassType) (*ast.ClassMemberSection, error) {
	defer p.TraceMethod("Parser.ParseClassMemberSection")()

	res := &ast.ClassMemberSection{}
	classType.Members = append(classType.Members, res)

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
			res.Visibility = ast.CvPublic
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
	if len(fieldList) > 0 {
		res.ClassFieldList = fieldList
	}

	methodList, err := p.ParseClassMethodList()
	if err != nil {
		return nil, err
	}
	if len(methodList) > 0 {
		res.ClassMethodList = methodList
	}

	propList, err := p.ParseClassPropertyList(classType)
	if err != nil {
		return nil, err
	}
	if len(propList) > 0 {
		res.ClassPropertyList = propList
	}

	return res, nil
}

var (
	visibilityBreak = token.Some(
		token.UpperCase("PRIVATE"),
		token.UpperCase("PROTECTED"),
		token.UpperCase("PUBLIC"),
		token.UpperCase("PUBLISHED"),
	)

	propertyBreak = token.Some(
		visibilityBreak,
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
	p.NextToken()

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
	if err := p.Until(methodBreak, nil, func() error {
		if methodBreak.Predicate(p.CurrentToken()) {
			return QuitUntil
		}
		method, err := p.ParseClassMethod()
		if err != nil {
			return err
		}
		res = append(res, method)
		if p.CurrentToken().Is(token.Symbol(';')) {
			p.NextToken()
		}
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

	defer p.context.StackDeclMap()()

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
	if len(directiveList) > 0 {
		res.Directives = directiveList
	}

	return res, nil
}

func (p *Parser) ClassMethodDirectiveList() (ast.ClassMethodDirectiveList, error) {
	if p.CurrentToken().Is(token.Symbol(';')) {
		p.NextToken()
	}

	defer p.TraceMethod("Parser.ClassMethodDirectiveList")()

	res := ast.ClassMethodDirectiveList{}
	if err := p.Until(methodBreak, token.Symbol(';'), func() error {
		w := strings.ToUpper(p.CurrentToken().Value())
		if ast.ClassMethodDirectives.Include(w) {
			res = append(res, ast.ClassMethodDirective(w))
			p.NextToken()
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
	defer p.TraceMethod("Parser.ParseConstructorHeading")()

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
	defer p.TraceMethod("Parser.ParseDestructorHeading")()

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
	defer p.TraceMethod("Parser.ParseClassPropertyList")()

	res := ast.ClassPropertyList{}
	if err := p.Until(propertyBreak, nil, func() error {
		if propertyBreak.Predicate(p.CurrentToken()) {
			return QuitUntil
		}
		prop, err := p.ParseClassProperty(classType)
		if err != nil {
			return err
		}
		res = append(res, prop)
		if p.CurrentToken().Is(token.Symbol(';')) {
			p.NextToken()
		}
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

	p.Logf("Parser.ParseClassProperty #01")

	t0 := p.NextToken()
	res.Ident = ast.NewIdent(t0)
	defer p.TraceMethod("Parser.ParseClassProperty")()

	p.NextToken()

	intf, err := p.ParsePropertyInterface()
	if err != nil {
		p.Logf("Parser.ParseClassProperty #02")
		return nil, err
	}
	res.Interface = intf

	p.Logf("Parser.ParseClassProperty #03")
	if strings.ToUpper(p.CurrentToken().Value()) == "INDEX" {
		p.Logf("Parser.ParseClassProperty #04")
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			p.Logf("Parser.ParseClassProperty #05")
			return nil, err
		}
		res.Index = expr
	}
	p.Logf("Parser.ParseClassProperty #06")

	//   [READ Ident]
	if strings.ToUpper(p.CurrentToken().Value()) == "READ" {
		p.Logf("Parser.ParseClassProperty #07")
		t := p.NextToken()
		if decl := classType.FindMemberDecl(t.Value()); decl != nil {
			res.Read = ast.NewIdentRef(ast.NewIdent(t), decl)
		} else {
			p.Logf("Parser.ParseClassProperty #08")
			return nil, p.TokenErrorf("unknown member %s", t)
		}
		p.NextToken()
	}
	p.Logf("Parser.ParseClassProperty #09")

	//   [WRITE Ident]
	if strings.ToUpper(p.CurrentToken().Value()) == "WRITE" {
		p.Logf("Parser.ParseClassProperty #10")
		t := p.NextToken()
		if decl := classType.FindMemberDecl(t.Value()); decl != nil {
			res.Write = ast.NewIdentRef(ast.NewIdent(t), decl)
		} else {
			p.Logf("Parser.ParseClassProperty #11")
			return nil, p.TokenErrorf("unknown member %s", t)
		}
		p.NextToken()
	}
	p.Logf("Parser.ParseClassProperty #12")

	//   [STORED (Ident | Constant)]
	if strings.ToUpper(p.CurrentToken().Value()) == "STORED" {
		p.Logf("Parser.ParseClassProperty #13")
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
			p.Logf("Parser.ParseClassProperty #14")
			return nil, p.TokenErrorf("unknown member %s", t)
		}
		p.NextToken()
	}
	p.Logf("Parser.ParseClassProperty #15")

	//   [(DEFAULT ConstExpr) | NODEFAULT]
	if strings.ToUpper(p.CurrentToken().Value()) == "DEFAULT" {
		p.Logf("Parser.ParseClassProperty #16")
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			p.Logf("Parser.ParseClassProperty #17")
			return nil, err
		}
		res.Default = &ast.PropertyDefaultSpecifier{Value: expr}
	} else if strings.ToUpper(p.CurrentToken().Value()) == "NODEFAULT" {
		p.Logf("Parser.ParseClassProperty #18")
		v := true
		res.Default = &ast.PropertyDefaultSpecifier{NoDefault: &v}
	}

	p.Logf("Parser.ParseClassProperty #19")

	//   [IMPLEMENTS TypeId]
	if strings.ToUpper(p.CurrentToken().Value()) == "IMPLEMENTS" {
		p.Logf("Parser.ParseClassProperty #20")
		typeId, err := p.ParseTypeId()
		if err != nil {
			p.Logf("Parser.ParseClassProperty #21")
			return nil, err
		}
		res.Implements = typeId
	}
	p.Logf("Parser.ParseClassProperty #22")

	//   [PortabilityDirective]
	//    TODO

	return res, nil
}

func (p *Parser) ParsePropertyInterface() (*ast.PropertyInterface, error) {
	res := &ast.PropertyInterface{}
	if p.CurrentToken().Is(token.Symbol('[')) {
		defer p.context.StackDeclMap()()

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
