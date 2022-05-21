unit Unit4;

interface

var
  Foo: string = 'Foo@Unit4';
  Bar: string = 'Bar@Unit4';

  Project1: string = 'Project1@Unit4';

implementation


initialization
  Writeln('-- Unit4 Initialization----');
  Writeln(Foo);
  Writeln(Bar);
  Writeln(Project1);

finalization
  Writeln('-- Unit4 Finalization----');
  Writeln(Foo);
  Writeln(Bar);
  Writeln(Project1);
  Readln;

end.
 