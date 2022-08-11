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
| :--: | ----------- | ----: |---------:|
|  🔖  | TODO        |    19 |  14.7% |
|  🚧  | In progress |     7 |   5.4% |
|  ✔️  | Done        |   103 | **79.8%** |
|     | Total        | 129 |  100.0%  |

See [Grammer.md](./Grammer.md) for more details.
