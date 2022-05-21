program simple_proj1;
{$APPTYPE CONSOLE}
uses
  SysUtils,
  call_inc in 'call_inc.pas',
  cnt in 'subdir1\cnt.pas';

begin
  CallInc;
  Readln;
end.