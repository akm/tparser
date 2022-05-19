unit cnt;

interface

var Count: Integer = 0;
procedure Inc;

implementation

procedure Inc;
begin
   Count := Count + 1;
end;

end.
