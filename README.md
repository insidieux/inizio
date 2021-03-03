# inizio

Golang project standard layout generator

[![GitHub Workflow Status](https://github.com/insidieux/inizio/workflows/Push/badge.svg)](https://github.com/insidieux/inizio/actions/workflows/push.yml?query=branch%3Amaster+event%3Apush)
[![Go Report Card](https://goreportcard.com/badge/github.com/insidieux/inizio)](https://goreportcard.com/report/github.com/insidieux/inizio)
[![codecov](https://codecov.io/gh/insidieux/inizio/branch/master/graph/badge.svg?token=BI6HEMPLB1)](https://codecov.io/gh/insidieux/inizio/branch/master)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/insidieux/inizio)

inizio is a simple binary, which allows generating/bootstrapping golang project with predefined layout. 
This project is easy can be extended, cause it also supports plugins for generation, based on [go-plugin](https://github.com/hashicorp/go-plugin) package. 

## Installing

Install `inizio` by running:

```shell
go get github.com/insidieux/inizio/cmd/inizio
```

Ensure that `$GOPATH/bin` is added to your `$PATH`.

## Documentation

- [User guide][]
- [Contributing guide][]
- [Write plugin guide][]

[User guide]: ./docs/user-guide.md
[Contributing guide]: ./docs/contributing.md
[Write plugin guide]: ./docs/write-plugin-guide.md


## License

[Apache][]

[Apache]: ./LICENSE
