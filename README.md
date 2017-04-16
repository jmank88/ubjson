# UBJSON [![GoDoc](https://godoc.org/github.com/jmank88/ubjson?status.svg)](https://godoc.org/github.com/jmank88/ubjson) [![Build Status](https://travis-ci.org/jmank88/ubjson.svg)](https://travis-ci.org/jmank88/ubjson) [![Go Report Card](https://goreportcard.com/badge/github.com/jmank88/ubjson)](https://goreportcard.com/report/github.com/jmank88/ubjson)
[Universal Binary JSON](http://ubjson.org/) library for Go.

# TODO
- Tag support: consider json, unless compelling reason for new ubjson
- Cache struct reflection on first run: fieldnames/tags/count

## Open Questions

- Should DecodeInt return int64 for machines where int is < 64bits
- Should Object{Encoder|Decoder} catch duplicate keys?
