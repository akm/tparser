unit Unit1;

interface

uses
  SysUtils, Unit4;

var
  Foo: string = 'Foo@Unit1';
  Bar: string = 'Bar@Unit1';

  Unit2: string = 'Unit2@Unit1';
  Project1: string = 'Project1@Unit1';

implementation

initialization
  Writeln('-- Unit1 Initialization----');
  Writeln(Foo);
  Writeln(Bar);
  Writeln(Unit2);

finalization
  Writeln('-- Unit1 Finalization----');
  Writeln(Foo);
  Writeln(Bar);
  Writeln(Unit2);
  Readln;

end.
