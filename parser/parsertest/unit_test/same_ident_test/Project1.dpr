program Project1;
{$APPTYPE CONSOLE}
uses
  SysUtils,
  Unit1 in 'Unit1.pas',
  Unit2 in 'Unit2.pas',
  Unit3 in 'Unit3.pas',
  Unit4 in 'Unit4.pas';

var
  // Unit1: string; // ErrorŽ¯•ÊŽq‚Ì‘½d’è‹`
  Bar: string = 'Bar@Project1';
  Baz: string = 'Baz@Project1';
begin
  Writeln('-- Project1 ----');
  Writeln(Foo); // => Foo@Unit4
  Writeln(Bar); // => Bar@Project1
  Writeln(Unit1.Unit2); // => Unit2@Unit1
  Writeln(Unit2.Unit1); // => Unit1@Unit2
  Writeln(Project1.Bar); // => Bar@Project1
  // Writeln(Project1.Unit4.Project1);
  Writeln(Project1.Baz); // => Baz@Project1
  Readln;
end.
