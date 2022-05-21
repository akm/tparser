unit bar;

interface

procedure Inc;
function Get: Integer;

implementation

var Count: Integer = 0;

procedure Inc;
begin
   Count := Count + 1; // Delphi allows semi-colon removed.
end;

function Get: Integer;
begin
   Result := Count;
end;

end.
