program example1;
{$APPTYPE CONSOLE}
uses
  SysUtils,
  foo in 'foo.pas',
  bar in 'subdir1\bar.pas';

begin
   foo.Process;
   Readln;
end.
