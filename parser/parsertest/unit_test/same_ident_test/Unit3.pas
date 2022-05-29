unit Unit3;

interface

uses
  Unit1, Unit2;


implementation

var
  Baz: string = 'Baz@Unit3';

initialization
  Writeln('-- Unit3 Initialization----');
  Writeln(Foo); // => Foo@Uni2
  Writeln(Bar); // => Foo@Unit2
  Writeln(Baz);
  Writeln(Unit3.Baz);
  Writeln(Project1); // => Project1@Unit2

finalization
  Writeln('-- Unit3 Finalization----');
  Writeln(Foo); // => Foo@Uni2
  Writeln(Bar); // => Foo@Unit2
  Writeln(Project1); // => Project1@Unit2
  Readln;

end.
