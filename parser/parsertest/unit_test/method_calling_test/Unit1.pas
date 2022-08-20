unit Unit1;

interface

type
  T0 = class
  public
    procedure Foo; virtual; abstract;
  end;

  T1 = class(T0)
  public
    procedure Foo; override;
  end;

  T2 = class(T1)
  public
    procedure Foo; override;
  end;

  T3 = class(T1)
  public
    procedure Foo; reintroduce;
  end;

  T4 = class(T1)
  public
    procedure Foo; // without directives like override or reintroduce;
  end;


implementation

procedure T1.Foo;
begin
  Writeln('T1.Foo');
end;

procedure T2.Foo;
begin
  Writeln('T2.Foo');
end;

procedure T3.Foo;
begin
  Writeln('T3.Foo');
end;

procedure T4.Foo;
begin
  Writeln('T4.Foo');
end;

end.
