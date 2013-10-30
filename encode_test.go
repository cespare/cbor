package cbor

import (
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
)

type testCase struct {
	input    interface{}
	expected string // hex bytes
}

// http://tools.ietf.org/html/rfc7049#appendix-A
var rfc7049TestCases = []testCase{
	// Positive integers
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

	// Negative integers
	{-1, "20"},
	{-10, "29"},
	{-100, "3863"},
	{-1000, "3903e7"},

	// Floats
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
	// TODO: missing float test cases

	// Simple objects
	{false, "f4"},
	{true, "f5"},
	{nil, "f6"},
	// TODO: undefined

	// TODO: unused simple numbers
	// TODO: all tagged examples

	// Byte strings
	{[]byte{}, "40"},
	{[]byte{1, 2, 3, 4}, "4401020304"},

	// Text string
	{"", "60"},
	{"a", "6161"},
	{"IETF", "6449455446"},
	{"\"\\", "62225c"},
	{"\u00fc", "62c3bc"},
	{"\u6c34", "63e6b0b4"},
	{"\U00010151", "64f0908591"},

	// Lists
	{[]int{}, "80"},
	{[]int{1, 2, 3}, "83010203"},
	{[]interface{}{1, []int{2, 3}, []int{4, 5}}, "8301820203820405"},
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		"98190102030405060708090a0b0c0d0e0f101112131415161718181819"},

	// Maps
	{map[int]int{}, "a0"},
	{map[int]int{1: 2, 3: 4}, "a201020304"},
	{map[string]interface{}{"a": 1, "b": []int{2, 3}}, "a26161016162820203"},
	{[]interface{}{"a", map[string]string{"b": "c"}}, "826161a161626163"},
	{map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E"},
		"a56161614161626142616361436164614461656145"},
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
	{[]string{"a", "b", "c"}, "83616161626163"},
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

type errTestCase struct {
	input            interface{}
	expectedErrRegex string
}

var errTestCases = []errTestCase{
	{string([]byte{0xff, 0xfe, 0xfd}), `string is not valid UTF-8`},
}

func TestEncodingErrors(t *testing.T) {
	for _, test := range errTestCases {
		_, err := Marshal(test.input)
		if err == nil {
			t.Error("Expected an non-nil error, but err was nil.")
			continue
		}
		r, err2 := regexp.Compile(test.expectedErrRegex)
		if err2 != nil {
			t.Error(err2)
			continue
		}
		if !r.MatchString(err.Error()) {
			t.Errorf("Expected error to match /%s/ but got '%s'", r, err)
		}
	}
}
