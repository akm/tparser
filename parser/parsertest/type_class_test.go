package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestClassType(t *testing.T) {

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
						Opt: &ast.FpoVar,
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
									// function Realloc(var NewCapacity: Longint): Pointer; virtual;
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
								Visibility: ast.CvPublished,
								ClassMethodList: ast.ClassMethodList{
									// destructor Destroy; override;
									{
										Heading:    &ast.DestructorHeading{Ident: asttest.NewIdent("Destroy")},
										Directives: ast.ClassMethodDirectives{ast.CmdOverride},
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
										Directives: ast.ClassMethodDirectives{ast.CmdOverride},
									},
									// function Write(const Buffer; Count: Longint): Longint; override;
									{
										Heading: &ast.FunctionHeading{Type: ast.FtFunction,
											Ident: asttest.NewIdent("Write"),
											FormalParameters: ast.FormalParameters{{
												Opt: &ast.FpoConst,
												Parameter: &ast.Parameter{
													IdentList: asttest.NewIdentList("Buffer", "Count"),
													Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Longint")},
												},
											}},
											ReturnType: asttest.NewOrdIdent("Longint"),
										},
										Directives: ast.ClassMethodDirectives{ast.CmdOverride},
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
						Parameter: &ast.Parameter{
							IdentList: asttest.NewIdentList("X", "Y"),
							Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
						},
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
						Opt: &ast.FpoConst,
						Parameter: &ast.Parameter{
							IdentList: asttest.NewIdentList("Name"),
							Type:      &ast.ParameterType{Type: asttest.NewStringType("string")},
						},
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
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Index"),
								Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
							},
						},
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Value"),
								Type:      &ast.ParameterType{Type: asttest.NewTypeId(asttest.NewIdent("TObject"))},
							},
						},
					},
				},
			}
			// procedure SetPixel(X, Y: Integer; Value: TColor);
			methodDeclSetPixel := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetPixel"),
					FormalParameters: ast.FormalParameters{
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("X", "Y"),
								Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
							},
						},
						{
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Value"),
								Type:      &ast.ParameterType{Type: asttest.NewTypeId(asttest.NewIdent("TColor"))},
							},
						},
					},
				},
			}
			// procedure SetValue(const Name, Value: string);
			methodDeclSetValue := &ast.ClassMethod{
				Heading: &ast.FunctionHeading{
					Type:  ast.FtProcedure,
					Ident: asttest.NewIdent("SetValue"),
					FormalParameters: ast.FormalParameters{
						{
							Opt: &ast.FpoConst,
							Parameter: &ast.Parameter{
								IdentList: asttest.NewIdentList("Name", "Value"),
								Type:      &ast.ParameterType{Type: asttest.NewStringType("string")},
							},
						},
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
								Visibility: ast.CvPblic,
								ClassPropertyList: ast.ClassPropertyList{
									// property Objects[Index: Integer]: TObject read GetObject write SetObject;
									{
										Ident: asttest.NewIdent("Objects"),
										Interface: &ast.PropertyInterface{
											Parameters: ast.FormalParameters{
												{
													Parameter: &ast.Parameter{
														IdentList: asttest.NewIdentList("Index"),
														Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
													},
												},
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
												{
													Parameter: &ast.Parameter{
														IdentList: asttest.NewIdentList("X", "Y"),
														Type:      &ast.ParameterType{Type: asttest.NewOrdIdent("Integer")},
													},
												},
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
												{
													Opt: &ast.FpoConst,
													Parameter: &ast.Parameter{
														IdentList: asttest.NewIdentList("Name"),
														Type:      &ast.ParameterType{Type: asttest.NewStringType("string")},
													},
												},
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
}
