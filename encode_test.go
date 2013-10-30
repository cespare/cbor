package cbor

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"testing"
)

type testCase struct {
	input    interface{}
	expected string // hex bytes
}

// http://tools.ietf.org/html/rfc7049#appendix-A
var rfc7049TestCases = []testCase{
	{0, "00"},
	{1, "01"},
	{10, "0a"},
	{23, "17"},
	{24, "1818"},
	{25, "1819"},
	{100, "1864"},
	{1000, "1903e8"},
	{1000000, "1a000f4240"},
	{1000000000000, "1b000000e8d4a51000"},

	// Bignums?
	//{18446744073709551615, "1bffffffffffffffff"},
	//{18446744073709551616, "c249010000000000000000"},
	//{-18446744073709551616, "3bffffffffffffffff"},
	//{-18446744073709551617, "c349010000000000000000"},

	{-1, "20"},
	{-10, "29"},
	{-100, "3863"},
	{-1000, "3903e7"},

	// Small float test cases omitted because this package does compact small floats down to half-precision.
	// See additionalTestCases.
	{1.1, "fb3ff199999999999a"},
	{100000.0, "fa47c35000"},
	{3.4028234663852886e+38, "fa7f7fffff"},
	{1.0e+300, "fb7e37e43c8800759c"},
	{-4.1, "fbc010666666666666"},
	{math.Inf(1), "fa7f800000"},
	//{math.NaN(), "fb7ff8000000000000"}, // TODO: NaN?
	{math.Inf(-1), "faff800000"},

	{false, "f4"},
	{true, "f5"},
	{nil, "f6"},

	{"", "60"},
	{"a", "6161"},
	{"IETF", "6449455446"},
}

var additionalTestCases = []testCase{
	// Floats
	{0.0, "fa00000000"},
	{-0.0, "fa00000000"},
	{1.0, "fa3f800000"},
	{1.5, "fa3fc00000"},
	{65504.0, "fa477fe000"},
	{5.960464477539063e-08, "fa33800000"},
	{0.00006103515625, "fa38800000"},
	{-4.0, "fac0800000"},
}

var typesTestCases = []testCase{
	{int(0), "00"},
	{int8(0), "00"},
	{int64(0), "00"},
	{int64(0), "00"},
	{uint(0), "00"},
	{uint8(0), "00"},
	{uint64(0), "00"},
	{uint(1), "01"},
	{uint(10), "0a"},
	{uint(23), "17"},
	{uint(24), "1818"},
	{uint(25), "1819"},
	{uint(100), "1864"},
	{uint(1000), "1903e8"},
	{uint(1000000), "1a000f4240"},
	{uint(1000000000000), "1b000000e8d4a51000"},
}

func TestEncoding(t *testing.T) {
	for _, suite := range [][]testCase{rfc7049TestCases, additionalTestCases, typesTestCases} {
		for _, test := range suite {
			b, err := Marshal(test.input)
			if err != nil {
				t.Error(err)
				continue
			}
			actual := hex.EncodeToString(b)
			if test.expected != actual {
				parts := []string{}
				for _, b2 := range b {
					parts = append(parts, fmt.Sprintf("%08b", b2))
				}
				fmt.Println(strings.Join(parts, " "))
				t.Errorf("Input: %#v, expected: 0x%s, actual: 0x%s", test.input, test.expected, actual)
			}
		}
	}
}
