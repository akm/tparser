unit Unit2;

interface

var
  Foo: string = 'Foo@Unit2';
  Bar: string = 'Bar@Unit2';
  Unit1: string = 'Unit1@Unit2';
  Project1: string = 'Project1@Unit2';

implementation

uses
  Unit1;

initialization
  Writeln('-- Unit2 Initialization----');
  Writeln(Foo);
  Writeln(Bar);
  // Writeln(Unit1);
  Writeln(Unit1.Foo);
  Writeln(Unit1.Bar);
  // Writeln(Unit1.Unit1.Foo);


finalization
  Writeln('-- Unit2 Finalization----');
  Writeln(Foo);
  Writeln(Bar);
  // Writeln(Unit1);
  Writeln(Unit1.Foo);
  Writeln(Unit1.Bar);
  Readln;

end.
