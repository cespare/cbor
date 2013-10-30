# cbor

Package cbor is a Go library for encoding and decoding CBOR (Concise Binary Object Representation) data. This
package implements CBOR as described in [RFC 7049](http://tools.ietf.org/html/rfc7049).

**Status:** WIP, not ready for consumption.

## TODO (or decide not to do)

* Tag for encoding empty strings as nulls
* Option to turn off unicode validity checking for strings for dat speed
* Option to allow (a/all) lists to be encoded with indefinite length (or some streaming API)
