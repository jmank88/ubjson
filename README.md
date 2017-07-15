# UBJSON [![GoDoc](https://godoc.org/github.com/jmank88/ubjson?status.svg)](https://godoc.org/github.com/jmank88/ubjson) [![Build Status](https://travis-ci.org/jmank88/ubjson.svg)](https://travis-ci.org/jmank88/ubjson) [![Go Report Card](https://goreportcard.com/badge/github.com/jmank88/ubjson)](https://goreportcard.com/report/github.com/jmank88/ubjson)

A Go package implementing encoding and decoding of [Universal Binary JSON](http://ubjson.org/) (spec 12).

## Usage

Most types can be automatically encoded through reflection with the Marshal
and Unmarshal functions. Encoders and Decoders additionally provide type
specific methods. Custom encodings can be defined by implementing the Value
interface.

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

See the [GoDoc](https://godoc.org/github.com/jmank88/ubjson) for more information.
