package cbor

const (
	// The first 6 major types have constants equal to their byte value
	typePosInt     byte = 0
	typeNegInt          = 1
	typeByteString      = 2
	typeTextString      = 3
	typeArray           = 4
	typeMap             = 5
	typeTag             = 6
	typeMajor7          = 7 // Overloaded for multiple types

	// TODO: Is this useful?
	typeUnassigned = iota

	// All the major type 7 types will conveniently have their value equal to their 5-bit ID value.
	typeFalse     = 20
	typeTrue      = 21
	typeNull      = 22
	typeUndefined = 23
	typeFloat16   = 25
	typeFloat32   = 26
	typeFloat64   = 27
	typeBreak     = 31
)

// Maps # bytes -> CBOR code
var additionalLength = [...]byte{
	1: 24,
	2: 25,
	4: 26,
	8: 27,
}
