package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/log/testlog"
)

func TestClassType(t *testing.T) {
	defer testlog.Setup(t)()

	RunTypeSection(t,
		"class example1",
		[]rune(`
type
  TMemoryStream = class(TCustomMemoryStream)
  private
    FCapacity: Longint;
    procedure SetCapacity(NewCapacity: Longint);
  protected
    function Realloc(var NewCapacity: Longint): Pointer; virtual;
    property Capacity: Longint read FCapacity write SetCapacity;
  public
    destructor Destroy; override;
    procedure Clear;
    procedure LoadFromStream(Stream: TStream);
    procedure LoadFromFile(const FileName: string);
    procedure SetSize(NewSize: Longint); override;
    function Write(const Buffer; Count: Longint): Longint; override;
  end;
`),
		func() ast.TypeSection {
			// FCapacity: Longint;
			fieldDeclFCapacity := &ast.ClassField{
				IdentList: asttest.NewIdentList("FCapacity"),
				Type:      asttest.NewOrdIdent("Longint"),
			}
			//  procedure SetCapacity(NewCapacity: Longint);
			methodDeclSetCapacity := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetCapacity"),
					FormalParameters: ast.FormalParameters{{
						Parameter: &ast.Parameter{
							IdentList: asttest.NewIdentList("NewCapacity"),
							Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Longint")},
						},
					}},
				},
			}

			return ast.TypeSection{
				&ast.TypeDecl{
					Ident: asttest.NewIdent("TMemoryStream"),
					Type: &ast.CustomClassType{
						Heritage: ast.ClassHeritage{asttest.NewTypeId("TCustomMemoryStream")},
						Members: ast.ClassMemberSections{
							&ast.ClassMemberSection{
								Visibility:      ast.CvPrivate,
								ClassFieldList:  ast.ClassFieldList{fieldDeclFCapacity},
								ClassMethodList: ast.ClassMethodList{methodDeclSetCapacity},
							},
							&ast.ClassMemberSection{
								Visibility: ast.CvProtected,
								ClassMethodList: ast.ClassMethodList{
									{
										Heading: &ast.FunctionHeading{
											Type:  ast.FtFunction,
											Ident: asttest.NewIdent("Realloc"),
											FormalParameters: ast.FormalParameters{{
												Opt: &ast.FpoVar,
												Parameter: &ast.Parameter{
													IdentList: asttest.NewIdentList("NewCapacity"),
													Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Longint")},
												},
											}},
											ReturnType: ast.NewPointerType(asttest.NewIdent("Pointer")),
										},
										Directives: ast.ClassMethodDirectiveList{ast.CmdVirtual},
									},
								},
								ClassPropertyList: ast.ClassPropertyList{
									// property Capacity: Longint read FCapacity write SetCapacity;
									{
										Ident:     asttest.NewIdent("Capacity"),
										Interface: &ast.PropertyInterface{Type: asttest.NewOrdIdent("Longint")},
										Read:      asttest.NewIdentRef("FCapacity", fieldDeclFCapacity.ToDeclarations()[0]),
										Write:     asttest.NewIdentRef("SetCapacity", methodDeclSetCapacity.ToDeclarations()[0]),
									},
								},
							},
							&ast.ClassMemberSection{
								Visibility: ast.CvPublic,
								ClassMethodList: ast.ClassMethodList{
									// destructor Destroy; override;
									{
										Heading:    &ast.DestructorHeading{Ident: asttest.NewIdent("Destroy")},
										Directives: ast.ClassMethodDirectiveList{ast.CmdOverride},
									},
									// procedure Clear;
									{Heading: &ast.FunctionHeading{Type: ast.FtProcedure, Ident: asttest.NewIdent("Clear")}},
									// procedure LoadFromStream(Stream: TStream);
									{
										Heading: &ast.FunctionHeading{Type: ast.FtProcedure,
											Ident: asttest.NewIdent("LoadFromStream"),
											FormalParameters: ast.FormalParameters{{
												Parameter: &ast.Parameter{
													IdentList: asttest.NewIdentList("Stream"),
													Type:      &ast.ParameterType{Type: asttest.NewTypeId("TStream")},
												},
											}},
										},
									},
									// procedure LoadFromFile(const FileName: string);
									{
										Heading: &ast.FunctionHeading{Type: ast.FtProcedure,
											Ident: asttest.NewIdent("LoadFromFile"),
											FormalParameters: ast.FormalParameters{{
												Opt: &ast.FpoConst,
												Parameter: &ast.Parameter{
													IdentList: asttest.NewIdentList("FileName"),
													Type:      &ast.ParameterType{Type: asttest.NewStringType("string")},
												},
											}},
										},
									},
									// procedure SetSize(NewSize: Longint); override;
									{
										Heading: &ast.FunctionHeading{Type: ast.FtProcedure,
											Ident: asttest.NewIdent("SetSize"),
											FormalParameters: ast.FormalParameters{{
												Parameter: &ast.Parameter{
													IdentList: asttest.NewIdentList("NewSize"),
													Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Longint")},
												},
											}},
										},
										Directives: ast.ClassMethodDirectiveList{ast.CmdOverride},
									},
									// function Write(const Buffer; Count: Longint): Longint; override;
									{
										Heading: &ast.FunctionHeading{Type: ast.FtFunction,
											Ident: asttest.NewIdent("Write"),
											FormalParameters: ast.FormalParameters{
												{
													Opt: &ast.FpoConst,
													Parameter: &ast.Parameter{
														IdentList: asttest.NewIdentList("Buffer"),
													},
												},
												{
													Parameter: &ast.Parameter{
														IdentList: asttest.NewIdentList("Count"),
														Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Longint")},
													},
												},
											},
											ReturnType: asttest.NewOrdIdent("Longint"),
										},
										Directives: ast.ClassMethodDirectiveList{ast.CmdOverride},
									},
								},
							},
						},
					},
				},
			}
		}(),
	)

	RunTypeSection(t,
		"simple inheritance",
		[]rune(`type TSomeControl = class(TControl);`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: asttest.NewIdent("TSomeControl"),
				Type: &ast.CustomClassType{
					Heritage: ast.ClassHeritage{asttest.NewTypeId("TControl")},
				},
			},
		},
	)

	RunTypeSection(t,
		"class compatibility",
		[]rune(`type
TFigure = class(TObject);
TRectangle = class(TFigure);
TSquare = class(TRectangle);
`),
		func() ast.TypeSection {
			classDeclFigure := &ast.TypeDecl{
				Ident: asttest.NewIdent("TFigure"),
				Type: &ast.CustomClassType{
					Heritage: ast.ClassHeritage{asttest.NewTypeId("TObject")},
				},
			}
			classDeclRectangle := &ast.TypeDecl{
				Ident: asttest.NewIdent("TRectangle"),
				Type: &ast.CustomClassType{
					Heritage: ast.ClassHeritage{asttest.NewTypeId("TFigure", classDeclFigure.ToDeclarations()[0])},
				},
			}
			classDeclSquare := &ast.TypeDecl{
				Ident: asttest.NewIdent("TSquare"),
				Type: &ast.CustomClassType{
					Heritage: ast.ClassHeritage{asttest.NewTypeId("TRectangle", classDeclRectangle.ToDeclarations()[0])},
				},
			}
			return ast.TypeSection{
				classDeclFigure,
				classDeclRectangle,
				classDeclSquare,
			}
		}(),
	)

	// The following type sections are invalid. The class declarated by forward declaration must be declarated completely in the same type section.
	// type
	// 		TFigure = class; // forward declaration
	// 		TDrawing = class
	// 			Figure: TFigure;
	// 		end;
	// type
	// 		TFigure = class // defining declaration
	// 			Drawing: TDrawing;
	// 		end;

	RunTypeSection(t,
		"Forward declarations and mutually dependent classes",
		[]rune(`
type
	TFigure = class; // forward declaration
	TDrawing = class
		Figure: TFigure;
	end;
	TFigure = class // defining declaration
		Drawing: TDrawing;
	end;
`),
		func() ast.TypeSection {
			classDeclFigure0 := &ast.TypeDecl{
				Ident: asttest.NewIdent("TFigure"),
				// Type:  &ast.ForwardDeclaredClassType{},
			}
			classDeclDrawing := &ast.TypeDecl{
				Ident: asttest.NewIdent("TDrawing"),
				Type: &ast.CustomClassType{
					Members: ast.ClassMemberSections{
						&ast.ClassMemberSection{
							Visibility: ast.CvPrivate,
							ClassFieldList: ast.ClassFieldList{
								&ast.ClassField{
									IdentList: asttest.NewIdentList("Figure"),
									Type:      asttest.NewTypeId("TFigure", classDeclFigure0.ToDeclarations()[0]),
								},
							},
						},
					},
				},
			}
			classFigure := &ast.CustomClassType{
				Members: ast.ClassMemberSections{
					&ast.ClassMemberSection{
						Visibility: ast.CvPrivate,
						ClassFieldList: ast.ClassFieldList{
							&ast.ClassField{
								IdentList: asttest.NewIdentList("Drawing"),
								Type:      asttest.NewTypeId("TDrawing", classDeclDrawing.ToDeclarations()[0]),
							},
						},
					},
				},
			}
			classDeclFigure1 := &ast.TypeDecl{
				Ident: asttest.NewIdent("Figure"),
				Type:  classFigure,
			}
			classDeclFigure0.Type = &ast.ForwardDeclaredClassType{
				Actual: classFigure,
			}
			return ast.TypeSection{
				classDeclFigure0,
				classDeclDrawing,
				classDeclFigure1,
			}
		}(),
	)

	RunTypeSection(t,
		"array properties",
		[]rune(`type
TArrayPropExample1 = class
private
	function GetObject(Index: Integer): TObject;
	function GetPixel(X, Y: Integer): TColor;
	function GetValue(const Name: string): string;
	procedure SetObject(Index: Integer; Value: TObject);
	procedure SetPixel(X, Y: Integer; Value: TColor);
	procedure SetValue(const Name, Value: string);
public
	property Objects[Index: Integer]: TObject read GetObject write SetObject;
	property Pixels[X, Y: Integer]: TColor read GetPixel write SetPixel;
	property Values[const Name: string]: string read GetValue write SetValue;
end;
`),
		func() ast.TypeSection {
			// function GetObject(Index: Integer): TObject;
			methodDeclGetObject := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: asttest.NewIdent("GetObject"),
					FormalParameters: ast.FormalParameters{{
						Parameter: &ast.Parameter{
							IdentList: asttest.NewIdentList("Index"),
							Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
						},
					}},
					ReturnType: asttest.NewTypeId(asttest.NewIdent("TObject")),
				},
			}
			// function GetPixel(X, Y: Integer): TColor;
			methodDeclGetPixel := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: asttest.NewIdent("GetPixel"),
					FormalParameters: ast.FormalParameters{{
						Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("X", "Y"), Type: &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")}},
					}},
					ReturnType: asttest.NewTypeId(asttest.NewIdent("TColor")),
				},
			}
			// function GetValue(const Name: string): string;
			methodDeclGetValue := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: asttest.NewIdent("GetValue"),
					FormalParameters: ast.FormalParameters{{
						Opt:       &ast.FpoConst,
						Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Name"), Type: &ast.ParameterType{Type: asttest.NewStringType("string")}},
					}},
					ReturnType: asttest.NewTypeId(asttest.NewIdent("string")),
				},
			}
			// procedure SetObject(Index: Integer; Value: TObject);
			methodDeclSetObject := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetObject"),
					FormalParameters: ast.FormalParameters{
						{Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Index"), Type: &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")}}},
						{Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Value"), Type: &ast.ParameterType{Type: asttest.NewTypeId(asttest.NewIdent("TObject"))}}},
					},
				},
			}
			// procedure SetPixel(X, Y: Integer; Value: TColor);
			methodDeclSetPixel := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetPixel"),
					FormalParameters: ast.FormalParameters{
						{Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("X", "Y"), Type: &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")}}},
						{Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Value"), Type: &ast.ParameterType{Type: asttest.NewTypeId(asttest.NewIdent("TColor"))}}},
					},
				},
			}
			// procedure SetValue(const Name, Value: string);
			methodDeclSetValue := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetValue"),
					FormalParameters: ast.FormalParameters{
						{Opt: &ast.FpoConst, Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Name", "Value"), Type: &ast.ParameterType{Type: asttest.NewStringType("string")}}},
					},
				},
			}

			return ast.TypeSection{
				&ast.TypeDecl{
					Ident: asttest.NewIdent("TArrayPropExample1"),
					Type: &ast.CustomClassType{
						Members: ast.ClassMemberSections{
							&ast.ClassMemberSection{
								Visibility: ast.CvPrivate,
								ClassMethodList: ast.ClassMethodList{
									methodDeclGetObject,
									methodDeclGetPixel,
									methodDeclGetValue,
									methodDeclSetObject,
									methodDeclSetPixel,
									methodDeclSetValue,
								},
							},
							&ast.ClassMemberSection{
								Visibility: ast.CvPublic,
								ClassPropertyList: ast.ClassPropertyList{
									// property Objects[Index: Integer]: TObject read GetObject write SetObject;
									{
										Ident: asttest.NewIdent("Objects"),
										Interface: &ast.PropertyInterface{
											Parameters: ast.FormalParameters{
												{Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Index"), Type: &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")}}},
											},
										},
										Read:  asttest.NewIdentRef("GetObject", methodDeclGetObject.ToDeclarations()[0]),
										Write: asttest.NewIdentRef("SetObject", methodDeclSetObject.ToDeclarations()[0]),
									},
									// property Pixels[X, Y: Integer]: TColor read GetPixel write SetPixel;
									{
										Ident: asttest.NewIdent("Pixels"),
										Interface: &ast.PropertyInterface{
											Parameters: ast.FormalParameters{
												{Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("X", "Y"), Type: &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")}}},
											},
										},
										Read:  asttest.NewIdentRef("GetPixel", methodDeclGetPixel.ToDeclarations()[0]),
										Write: asttest.NewIdentRef("SetPixel", methodDeclSetPixel.ToDeclarations()[0]),
									},
									// property Values[const Name: string]: string read GetValue write SetValue;
									{
										Ident: asttest.NewIdent("Values"),
										Interface: &ast.PropertyInterface{
											Parameters: ast.FormalParameters{
												{Opt: &ast.FpoConst, Parameter: &ast.Parameter{IdentList: asttest.NewIdentList("Name"), Type: &ast.ParameterType{Type: asttest.NewStringType("string")}}},
											},
										},
										Read:  asttest.NewIdentRef("GetValue", methodDeclGetValue.ToDeclarations()[0]),
										Write: asttest.NewIdentRef("SetValue", methodDeclSetValue.ToDeclarations()[0]),
									},
								},
							},
						},
					},
				},
			}
		}(),
	)

	RunTypeSection(t,
		"array properties",
		[]rune(`
type
	TRectangle = class
	private
		FCoordinates: array[0..3] of Longint;
		function GetCoordinate(Index: Integer): Longint;
		procedure SetCoordinate(Index: Integer; Value: Longint);
	public
		property Left: Longint index 0 read GetCoordinate write SetCoordinate;
		property Top: Longint index 1 read GetCoordinate write SetCoordinate;
		property Right: Longint index 2 read GetCoordinate write SetCoordinate;
		property Bottom: Longint index 3 read GetCoordinate write SetCoordinate;
		property Coordinates[Index: Integer]: Longint read GetCoordinate write SetCoordinate;
	end;
`),
		func() ast.TypeSection {
			// FCoordinates: array[0..3] of Longint;
			fieldDeclFCoordinates := &ast.ClassField{
				IdentList: asttest.NewIdentList("FCoordinates"),
				Type: &ast.ArrayType{
					IndexTypes: []ast.OrdinalType{
						&ast.SubrangeType{
							Low:  asttest.NewConstExpr(asttest.NewNumber("0")),
							High: asttest.NewConstExpr(asttest.NewNumber("3")),
						},
					},
					BaseType: asttest.NewOrdIdent("Longint"),
				},
			}
			// function GetCoordinate(Index: Integer): Longint;
			methodDeclGetCoordinate := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtFunction,
					Ident: asttest.NewIdent("GetCoordinate"),
					FormalParameters: ast.FormalParameters{{
						Parameter: &ast.Parameter{
							IdentList: asttest.NewIdentList("Index"),
							Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
						},
					}},
					ReturnType: asttest.NewOrdIdent("Longint"),
				},
			}
			// procedure SetCoordinate(Index: Integer; Value: Longint);
			methodDeclSetCoordinate := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetCoordinate"),
					FormalParameters: ast.FormalParameters{
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Index"),
								Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
							},
						},
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Value"),
								Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Longint")},
							},
						},
					},
				},
			}
			declGetCoordinate := methodDeclGetCoordinate.ToDeclarations()[0]
			declSetCoordinate := methodDeclSetCoordinate.ToDeclarations()[0]

			return ast.TypeSection{
				&ast.TypeDecl{
					Ident: asttest.NewIdent("TArrayPropExample1"),
					Type: &ast.CustomClassType{
						Members: ast.ClassMemberSections{
							&ast.ClassMemberSection{
								Visibility: ast.CvPrivate,
								ClassFieldList: ast.ClassFieldList{
									fieldDeclFCoordinates,
								},
								ClassMethodList: ast.ClassMethodList{
									methodDeclGetCoordinate,
									methodDeclSetCoordinate,
								},
							},
							&ast.ClassMemberSection{
								Visibility: ast.CvPublic,
								ClassPropertyList: ast.ClassPropertyList{
									// property Left: Longint index 0 read GetCoordinate write SetCoordinate;
									{
										Ident:     asttest.NewIdent("Left"),
										Interface: &ast.PropertyInterface{Type: asttest.NewOrdIdent("Longint")},
										Index:     asttest.NewExpression(asttest.NewNumber("0")),
										Read:      asttest.NewIdentRef("GetObject", declGetCoordinate),
										Write:     asttest.NewIdentRef("SetObject", declSetCoordinate),
									},
									// property Top: Longint index 1 read GetCoordinate write SetCoordinate;
									{
										Ident:     asttest.NewIdent("Top"),
										Interface: &ast.PropertyInterface{Type: asttest.NewOrdIdent("Longint")},
										Index:     asttest.NewExpression(asttest.NewNumber("1")),
										Read:      asttest.NewIdentRef("GetObject", declGetCoordinate),
										Write:     asttest.NewIdentRef("SetObject", declSetCoordinate),
									},
									// property Right: Longint index 2 read GetCoordinate write SetCoordinate;
									{
										Ident:     asttest.NewIdent("Right"),
										Interface: &ast.PropertyInterface{Type: asttest.NewOrdIdent("Longint")},
										Index:     asttest.NewExpression(asttest.NewNumber("2")),
										Read:      asttest.NewIdentRef("GetObject", declGetCoordinate),
										Write:     asttest.NewIdentRef("SetObject", declSetCoordinate),
									},
									// property Bottom: Longint index 3 read GetCoordinate write SetCoordinate;
									{
										Ident:     asttest.NewIdent("Bottom"),
										Interface: &ast.PropertyInterface{Type: asttest.NewOrdIdent("Longint")},
										Index:     asttest.NewExpression(asttest.NewNumber("3")),
										Read:      asttest.NewIdentRef("GetObject", declGetCoordinate),
										Write:     asttest.NewIdentRef("SetObject", declSetCoordinate),
									},
									// property Coordinates[Index: Integer]: Longint read GetCoordinate write SetCoordinate;
									{
										Ident: asttest.NewIdent("Coordinates"),
										Interface: &ast.PropertyInterface{
											Parameters: ast.FormalParameters{
												{Parameter: &ast.Parameter{
													IdentList: asttest.NewIdentList("Index"),
													Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
												}},
											},
											Type: asttest.NewOrdIdent("Longint"),
										},
										Read:  asttest.NewIdentRef("GetObject", declGetCoordinate),
										Write: asttest.NewIdentRef("SetObject", declSetCoordinate),
									},
								},
							},
						},
					},
				},
			}
		}(),
	)

	RunTypeSection(t,
		"Property overrides and redeclarations",
		[]rune(`
type
	TAncestor = class
	protected
		FSize: Integer;
	private
		function GetText: string;
		procedure SetText(const Value: string);
	private
		FColor: TColor;
		procedure SetColor(Value: TColor);
	protected
		property Size: Integer read FSize;
		property Text: string read GetText write SetText;
		property Color: TColor read FColor write SetColor stored False;
	end;

	TDerived = class(TAncestor)
	private
		procedure SetSize(const Value: Integer);
	protected
		property Size write SetSize;
	published
		property Text;
		property Color stored True default clBlue;
	end;
`),
		func() ast.TypeSection {
			// FSize: Integer;
			fieldDeclFSize := &ast.ClassField{IdentList: asttest.NewIdentList("FSize"), Type: asttest.NewOrdIdent("Integer")}
			// function GetText: string;
			methodDeclGetText := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{Type: ast.FtFunction, Ident: asttest.NewIdent("GetText"), ReturnType: asttest.NewStringType("string")},
			}
			// procedure SetText(const Value: string);
			methodDeclSetText := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetText"),
					FormalParameters: ast.FormalParameters{
						{
							Opt: &ast.FpoConst,
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Value"), Type: &ast.ParameterType{Type: asttest.NewStringType("string")}},
						},
					},
				},
			}

			// FColor: TColor;
			fieldDeclFColor := &ast.ClassField{IdentList: asttest.NewIdentList("FColor"), Type: asttest.NewTypeId("TColor")}
			// procedure SetColor(Value: TColor);
			methodDeclSetColor := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetColor"),
					FormalParameters: ast.FormalParameters{
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Value"), Type: &ast.ParameterType{Type: asttest.NewTypeId("TColor")}},
						},
					},
				},
			}

			// property Size: Integer read FSize;
			propertyDeclSize1 := &ast.ClassProperty{
				Ident:     asttest.NewIdent("Size"),
				Interface: &ast.PropertyInterface{Type: asttest.NewOrdIdent("Integer")},
				Read:      asttest.NewIdentRef("FSize", fieldDeclFSize.ToDeclarations()[0]),
			}
			// property Text: string read GetText write SetText;
			propertyDeclText1 := &ast.ClassProperty{
				Ident:     asttest.NewIdent("Text"),
				Interface: &ast.PropertyInterface{Type: asttest.NewStringType("string")},
				Read:      asttest.NewIdentRef("GetText", methodDeclGetText.ToDeclarations()[0]),
				Write:     asttest.NewIdentRef("SetText", methodDeclSetText.ToDeclarations()[0]),
			}
			// property Color: TColor read FColor write SetColor stored False;
			propertyDeclColor1 := &ast.ClassProperty{
				Ident:     asttest.NewIdent("Color"),
				Interface: &ast.PropertyInterface{Type: asttest.NewTypeId("TColor")},
				Read:      asttest.NewIdentRef("FColor", fieldDeclFColor.ToDeclarations()[0]),
				Write:     asttest.NewIdentRef("SetColor", methodDeclSetColor.ToDeclarations()[0]),
				Stored:    &ast.PropertyStoredSpecifier{Constant: ext.BoolPtr(true)},
			}

			classDeclAncestor := &ast.TypeDecl{
				Ident: asttest.NewIdent("TAncestor"),
				Type: &ast.CustomClassType{
					Members: ast.ClassMemberSections{
						&ast.ClassMemberSection{
							Visibility:     ast.CvProtected,
							ClassFieldList: ast.ClassFieldList{fieldDeclFSize},
						},
						&ast.ClassMemberSection{
							Visibility:      ast.CvPrivate,
							ClassMethodList: ast.ClassMethodList{methodDeclGetText, methodDeclSetText},
						},
						&ast.ClassMemberSection{
							Visibility:      ast.CvPrivate,
							ClassFieldList:  ast.ClassFieldList{fieldDeclFColor},
							ClassMethodList: ast.ClassMethodList{methodDeclSetColor},
						},
						&ast.ClassMemberSection{
							Visibility:        ast.CvProtected,
							ClassPropertyList: ast.ClassPropertyList{propertyDeclSize1, propertyDeclText1, propertyDeclColor1},
						},
					},
				},
			}

			// 	procedure SetSize(const Value: Integer);
			methodDeclSetSize := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetSize"),
					FormalParameters: ast.FormalParameters{
						{
							Opt: &ast.FpoConst,
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Value"), Type: &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")}},
						},
					},
				},
			}

			// 	property Size write SetSize;
			propertyDeclSize2 := &ast.ClassProperty{
				Ident:  asttest.NewIdent("Size"),
				Write:  asttest.NewIdentRef("SetSize", methodDeclSetSize.ToDeclarations()[0]),
				Parent: propertyDeclSize1,
			}
			// 	property Text;
			propertyDeclText2 := &ast.ClassProperty{
				Ident:  asttest.NewIdent("Text"),
				Parent: propertyDeclText1,
			}
			// 	property Color stored True default clBlue;
			propertyDeclColor2 := &ast.ClassProperty{
				Ident:     asttest.NewIdent("Color"),
				Interface: &ast.PropertyInterface{Type: asttest.NewTypeId("TColor")},
				Stored:    &ast.PropertyStoredSpecifier{Constant: ext.BoolPtr(true)},
				Default:   &ast.PropertyDefaultSpecifier{Value: asttest.NewConstExpr(ast.NewValueFactor("clBlue"))},
			}

			classDeclDerived := &ast.TypeDecl{
				Ident: asttest.NewIdent("TDerived"),
				Type: &ast.CustomClassType{
					Heritage: ast.ClassHeritage{
						asttest.NewTypeId("TAncestor", classDeclAncestor.ToDeclarations()[0]),
					},
					Members: ast.ClassMemberSections{
						&ast.ClassMemberSection{
							Visibility:      ast.CvPrivate,
							ClassMethodList: ast.ClassMethodList{methodDeclSetSize},
						},
						&ast.ClassMemberSection{
							Visibility:        ast.CvProtected,
							ClassPropertyList: ast.ClassPropertyList{propertyDeclSize2},
						},
						&ast.ClassMemberSection{
							Visibility:        ast.CvPrivate,
							ClassPropertyList: ast.ClassPropertyList{propertyDeclText2, propertyDeclColor2},
						},
					},
				},
			}

			return ast.TypeSection{classDeclAncestor, classDeclDerived}
		}(),
	)
}
