# go-test-color

[![Go](https://github.com/haunt98/go-test-color/workflows/Go/badge.svg?branch=main)](https://github.com/actions/setup-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/haunt98/go-test-color.svg)](https://pkg.go.dev/github.com/haunt98/go-test-color)
[![Latest Version](https://img.shields.io/github/v/tag/haunt98/go-test-color)](https://github.com/haunt98/go-test-color/tags)

Run `go test` with color.

## Install

With Go version `>= 1.16`:

```sh
go install github.com/haunt98/go-test-color@latest
```

With Go version `< 1.16`:

```sh
GO111module=on go get -u github.com/haunt98/go-test-color
```

## Usage

```sh
# Simply replace go test with go-test-color
go-test-color -v ./...
```

## Thanks

- [fatih/color](https://github.com/fatih/color)
- [rakyll/gotest](https://github.com/rakyll/gotest)
