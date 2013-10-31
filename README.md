# cbor

Package cbor is a Go library for encoding and decoding CBOR (Concise Binary Object Representation) data. This
package implements CBOR as described in [RFC 7049](http://tools.ietf.org/html/rfc7049).

**Status:** WIP, not ready for consumption.

## TODO (or decide not to do)

* Tag for encoding empty strings as nulls
* Option to turn off unicode validity checking for strings for dat speed
* Option to allow (a/all) lists to be encoded with indefinite length (or some streaming API)
* A lower-level API to help people write arbitrary CBOR messages (or at least help them when writing their own
  `Marshaler`/`Unmarshaler` implementations).
* Channel, complex, and function values cannot be marshaled by encoding/json, and I'm following suit here. We
  might be able to put complex numbers into a byte string with a tag or something (but it's not a predefined
  tag, so maybe don't bother).
* Handle anonymous struct fields the same way encoding/json does. Skipping this for now because it rachets up
  the complexity significantly.
* Better test case coverage for error cases in the encoder.
* encoding/json.Unmarshal will allow for type errors as it decodes and still give the user a best-effort
  decoded value as well as the error. Is this worth doing?
* json.{En,De}coder equivalents.
