unit foo;

interface

procedure Process;

implementation

uses SysUtils, bar;

procedure Process;
begin
   Inc;
   Writeln( bar.Get );
   bar.Inc;
   Writeln( bar.Get );
end;

end.
