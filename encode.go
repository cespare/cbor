package cbor

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"unicode/utf8"
)

func Marshal(v interface{}) ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

func (e *encodeState) error(err error) {
	panic(err)
}

type Marshaler interface {
	MarshalCBOR() ([]byte, error)
}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("cbor: unsupported type: %s", e.Type)
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return fmt.Sprintf("cbor: unsupported value: %s", e.Str)
}

type MarshalerError struct {
	Type reflect.Type
	Err  error
}

type InvalidUTF8Error struct {
	Str string
}

func (e *InvalidUTF8Error) Error() string {
	return fmt.Sprintf("cbor: string is not valid UTF-8: %s", e.Str)
}

func (e *MarshalerError) Error() string {
	return fmt.Sprintf("cbor: error calling MarshalCBOR for type %s: %s", e.Type, e.Err)
}

func (e *encodeState) reflectValue(v reflect.Value) {
	if !v.IsValid() {
		e.writeSimple(typeNull)
		return
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		// T isn't a Marshaler. Check *T as well.
		if v.Kind() != reflect.Ptr && v.CanAddr() {
			m, ok = v.Addr().Interface().(Marshaler)
			if ok {
				v = v.Addr()
			}
		}
	}
	if ok && (v.Kind() != reflect.Ptr && !v.IsNil()) {
		b, err := m.MarshalCBOR()
		if err != nil {
			// TODO: encoding/json parses the output of MarshalJSON here to check its validity. Do we want to do
			// that? (Punt until after a reasonable decoder is written, anyway.)
			e.Write(b)
			return
		}
		e.error(&MarshalerError{v.Type(), err})
	}

	switch v.Kind() {
	case reflect.Bool:
		x := v.Bool()
		if x {
			e.writeSimple(typeTrue)
		} else {
			e.writeSimple(typeFalse)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := v.Int()
		typ := typePosInt
		negative := i < 0
		if negative {
			i = -1 - i
			typ = typeNegInt
		}
		e.writeMajorWithNumber(typ, uint64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.writeMajorWithNumber(typePosInt, v.Uint())

	// TODO: Float canonicalization?
	case reflect.Float32:
		e.WriteByte(makeIDByte(typeMajor7, additionalLength[4]))
		e.putUint32(math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		f := v.Float()
		f32 := float32(f)
		// See if f is representable as a float32.
		if f == float64(f32) {
			e.WriteByte(makeIDByte(typeMajor7, additionalLength[4]))
			e.putUint32(math.Float32bits(f32))
		} else {
			e.WriteByte(makeIDByte(typeMajor7, additionalLength[8]))
			e.putUint64(math.Float64bits(v.Float()))
		}
	case reflect.String:
		s := v.String()
		if !utf8.ValidString(s) {
			e.error(&InvalidUTF8Error{s})
		}
		e.writeMajorWithNumber(typeTextString, uint64(len(s)))
		e.WriteString(s)
	default:
		e.error(&UnsupportedTypeError{v.Type()})
	}
}

type encodeState struct {
	bytes.Buffer
}

// makeIDByte returns a byte with the top 3 bits set to the value of major (should be < 8) and the bottom 5
// bits set to value (should be < 32).
func makeIDByte(major, value byte) byte {
	// 0x1F = 0b0001_1111
	return (value & 0x1F) | (major << 5)
}

func (e *encodeState) writeSimple(typ byte) {
	switch typ {
	case typeFalse, typeTrue, typeNull, typeUndefined, typeBreak:
		e.WriteByte(makeIDByte(7, typ))
	default:
		panic("not a simple type")
	}
}

func (e *encodeState) putUint8(i uint8) {
	e.WriteByte(byte(i))
}

func (e *encodeState) putUint16(i uint16) {
	e.WriteByte(byte(i >> 8))
	e.WriteByte(byte(i))
}

func (e *encodeState) putUint32(i uint32) {
	e.WriteByte(byte(i >> 24))
	e.WriteByte(byte(i >> 16))
	e.WriteByte(byte(i >> 8))
	e.WriteByte(byte(i))
}

func (e *encodeState) putUint64(i uint64) {
	e.WriteByte(byte(i >> 56))
	e.WriteByte(byte(i >> 48))
	e.WriteByte(byte(i >> 40))
	e.WriteByte(byte(i >> 32))
	e.WriteByte(byte(i >> 24))
	e.WriteByte(byte(i >> 16))
	e.WriteByte(byte(i >> 8))
	e.WriteByte(byte(i))
}

// writeMajorWithNumber writes in the given major type and a count, encoded using CBOR's number encoding
// method where count < 24 is written in the last 5 bytes, < 256 are written with 1 extra byte, etc. This
// is used for number encoding as well as the lengths of arrays and maps.
func (e *encodeState) writeMajorWithNumber(major byte, count uint64) {
	// Canonically, numbers are put into the smallest possible representation.
	switch {
	case count < 24:
		e.WriteByte(makeIDByte(major, byte(count)))
	case count < 256:
		e.WriteByte(makeIDByte(major, additionalLength[1]))
		e.putUint8(uint8(count))
	case count < 65536:
		e.WriteByte(makeIDByte(major, additionalLength[2]))
		e.putUint16(uint16(count))
	case count < 4294967296:
		e.WriteByte(makeIDByte(major, additionalLength[4]))
		e.putUint32(uint32(count))
	default:
		e.WriteByte(makeIDByte(major, additionalLength[8]))
		e.putUint64(uint64(count))
	}
}

func (e *encodeState) marshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	e.reflectValue(reflect.ValueOf(v))
	return nil
}
