# UBJSON [![GoDoc](https://godoc.org/github.com/jmank88/ubjson?status.svg)](https://godoc.org/github.com/jmank88/ubjson) [![Build Status](https://github.com/jmank88/ubjson/workflows/Go%20Build%20and%20Test/badge.svg)](https://github.com/jmank88/ubjson/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/jmank88/ubjson)](https://goreportcard.com/report/github.com/jmank88/ubjson)

A Go package implementing encoding and decoding of [Universal Binary JSON](http://ubjson.org/) (spec 12).

## Features

- Type specific methods for built-in types.

- Automatic encoding via reflection for most types.

- Custom encoding via Value interface.

- Streaming support via Encoder/Decoder.

- Support for [optimized format](http://ubjson.org/type-reference/container-types/#optimized-format).

- Block format.

## Usage

```go
b, _ := ubjson.MarshalBlock(8)
// [U][8]

b, _ = ubjson.MarshalBlock("hello")
// [S][U][5][hello]

var v interface{}
...
b, _ = ubjson.Marshal(v)
// ...
```

See the [GoDoc](https://godoc.org/github.com/jmank88/ubjson) for more
information and [examples](https://godoc.org/github.com/jmank88/ubjson#pkg-examples).
