# Write plugin guide

## Bootstrap project

```shell
inizio /path/to/project/directory
cd /path/to/project/directory
```

## Write generator

To write new generator plugin, you must implement `generator.Generator` interface in your code and call predefined `plugin.Serve` with written implementation.

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

## Get dependency

```shell
go get github.com/insidieux/inizio/pkg/sdk/generator
go get github.com/insidieux/inizio/pkg/sdk/plugin
```

or

```shell
make vendor
```

## Build

```shell
go build \
  -mod vendor \
  -o /path/to/project/directory/bin/your-plugin-name \
  -v /path/to/project/directory/cmd/your-plugin-name
```

or 

```shell
make build
```

That's all. Your first plugin is ready.
