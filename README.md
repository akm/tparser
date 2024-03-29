# TParser

## Overview

TParser is Object Pascal source code parser for alanysis instead of compiling. It builds AST nodes from source code at this time. In the near future, it AST nodes is going to be serialized to RDB or something.
TParser is based on Delphi's [Object Pascal Guide](https://docs.embarcadero.com/products/rad_studio/cbuilder6/EN/CB6_ObjPascalLangGuide_EN.pdf) in [RAD Studio documents](https://docs.embarcadero.com/products/rad_studio/).
See [Grammer.md](./Grammer.md).

## Run tests

```
go test ./...
```

## Status

| Mark | State       | Count | Percentage |
| :--: | ----------- | ----: | ---------: |
|  🔖  | TODO        |    18 |      14.0% |
|  🚧  | In progress |     6 |       4.7% |
|  ✔️  | Done        |   105 |  **81.4%** |
|      | Total       |   129 |     100.0% |

See [Grammer.md](./Grammer.md) for more details.
