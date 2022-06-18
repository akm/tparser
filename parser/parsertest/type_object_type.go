package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestObjectType(t *testing.T) {

	RunUnitTest(t,
		"object with visibility, class method and inheritance with interface type",
		[]rune(`UNIT U1;
interface
implementation

type
  TFoo = object
  private
    FBaz: Integer;
  public
    class function Bar: boolean;
    property Baz: Integer read FBaz;
  end;

class function TFoo.Bar: boolean;
begin
  Result := True;
end;


type
  IFoo = interface
    procedure Hoge;
  end;

type
  TBar = object(TFoo, IFoo)
  protected
    function QueryInterface(const IID: TGUID; out Obj): HResult; stdcall;
    function _AddRef: Integer; stdcall;
    function _Release: Integer; stdcall;
  public
    procedure Hoge;
  end;

procedure TBar.Hoge;
begin
end;

function TBar.QueryInterface(const IID: TGUID; out Obj): HResult; stdcall;
begin
  Result := 1;
end;

function TBar._AddRef: Integer; stdcall;
begin
  Result := 1;
end;

function TBar._Release: Integer; stdcall;
begin
  Result := 1;
end;

end.
`),
		func() *ast.Unit {
			return &ast.Unit{
				Ident:                 asttest.NewIdent("U1"),
				InterfaceSection:      &ast.InterfaceSection{},
				ImplementationSection: &ast.ImplementationSection{},
			}
		}(),
	)

}
