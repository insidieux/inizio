# Write plugin guide

## Bootstrap project

```shell
inizio /path/to/project/directory
cd /path/to/project/directory
```

## Get dependency

```shell
go get github.com/insidieux/inizio/pkg/sdk/generator
go get github.com/insidieux/inizio/pkg/sdk/plugin
```

or 

```shell
make vendor
```

## Write generator

```go
package main

import (
	"context"

	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/insidieux/inizio/pkg/sdk/generator/plugin"
)

type (
	implementation struct{}
)

var (
	_ generator.Generator = &implementation{}
)

func (i *implementation) Run(ctx context.Context, options generator.RunOptions, values generator.RunValues) (generator.RunResult, error) {
	panic("Not implemented yet")
}

func main() {
	plugin.Serve(new(implementation))
}
```

## Build

```shell
go build
  -mod vendor
  -o /project/bin/your-plugin-name
  -v /project/cmd/your-plugin-name
```

or 

```shell
make build
```

That's all. Your first plugin is ready.
