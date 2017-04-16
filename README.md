# UBJSON
[Universal Binary JSON](http://ubjson.org/) for Go.

This is a work in progress with a partially function encoder, and a decoder shell.

Encode for easy use

Encode*Type* for faster use

# TODO
- Tag support: consider json, unless compelling reason for new ubjson
- Cache struct reflection on first run: fieldnames/tags/count

TODO example

//TODO godoc link
//TODO CI

## Open Questions

- Should DecodeInt return int64 for machines where int is < 64bits
- Should Object{Encoder|Decoder} catch duplicate keys?
