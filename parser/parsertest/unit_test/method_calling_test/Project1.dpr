program Project1;
{$APPTYPE CONSOLE}
uses
  SysUtils,
  Unit1 in 'Unit1.pas';

var
  t01, t02, t03, t04: T0;
  t11, t12, t13, t14: T1;
  t22: T2;
  t33: T3;
  t44: T4;
begin
  Writeln('T0');
  t01 := T1.Create;
  t02 := T2.Create;
  t03 := T3.Create;
  t04 := T4.Create;
  t01.Foo;
  t02.Foo;
  t03.Foo;
  t04.Foo;

  Writeln('T1');
  t11 := T1.Create;
  t12 := T2.Create;
  t13 := T3.Create;
  t14 := T4.Create;
  t11.Foo;
  t12.Foo;
  t13.Foo;
  t14.Foo;

  Writeln('T2');
  t22 := T2.Create;
  t22.Foo;

  Writeln('T3');
  t33 := T3.Create;
  t33.Foo;

  Writeln('T4');
  t44 := T4.Create;
  t44.Foo;

  Readln;
end.