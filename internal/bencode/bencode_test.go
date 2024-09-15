package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"simple dictionary", map[string]interface{}{"foo": "bar"}, "d3:foo3:bare"},
		{"complex dictionary", map[string]interface{}{"hello": 4, "amore": []any{"a", 22}}, "d5:amorel1:ai22ee5:helloi4ee"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Marshal(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDecodeString(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedValue  string
		expectedRemain string
		expectError    bool
	}{
		{"valid string", "4:spam", "spam", "", false},
		{"string with remaining", "2:hii100elove", "hi", "i100elove", false},
		{"empty string", "0:", "", "", false},
		{"invalid length", "4:sp", "", "", true},
		{"missing content", "10:", "", "", true},
		{"invalid format", "1x:ha", "", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, remain, err := decodeString(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
				assert.Equal(t, tc.expectedRemain, remain)
			}
		})
	}
}

func TestDecodeInteger(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedValue  int
		expectedRemain string
		expectError    bool
	}{
		{"positive integer", "i52e", 52, "", false},
		{"negative integer", "i-100ethisdoesnotmatter", -100, "thisdoesnotmatter", false},
		{"invalid format", "i35d2e", 0, "", true},
		{"missing end", "i2", 0, "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, remain, err := decodeInteger(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
				assert.Equal(t, tc.expectedRemain, remain)
			}
		})
	}
}

func TestDecodeList(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedValue  []interface{}
		expectedRemain string
		expectError    bool
	}{
		{"string list", "l4:spam4:eggse", []interface{}{"spam", "eggs"}, "", false},
		{"mixed list", "l4:spami52ee", []interface{}{"spam", 52}, "", false},
		{"empty list", "le", []interface{}{}, "", false},
		{"invalid list", "l4:", nil, "", true},
		{"unclosed list", "lhi", nil, "", true},
		{"incomplete list", "l2:hi", nil, "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, remain, err := decodeList(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
				assert.Equal(t, tc.expectedRemain, remain)
			}
		})
	}
}

func TestDecodeDictionary(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedValue  map[string]interface{}
		expectedRemain string
		expectError    bool
	}{
		{"empty dictionary", "de", map[string]interface{}{}, "", false},
		{"simple dictionary", "d3:foo3:bar5:helloi52ee", map[string]interface{}{"foo": "bar", "hello": 52}, "", false},
		{"invalid dictionary", "d", nil, "", true},
		{"missing value", "d5:hello", nil, "", true},
		{"invalid key", "diloee", nil, "", true},
		{"incomplete value", "d2:hii12df", nil, "", true},
		{"non-string key", "di33e2:hie", nil, "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, remain, err := decodeDictionary(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
				assert.Equal(t, tc.expectedRemain, remain)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedValue interface{}
		expectError   bool
	}{
		{"string", "5:hello", "hello", false},
		{"integer", "i52e", 52, false},
		{"list", "l5:helloi52ee", []interface{}{"hello", 52}, false},
		{"dictionary", "d3:foo3:bar5:helloi52ee", map[string]interface{}{"foo": "bar", "hello": 52}, false},
		{"invalid input", "x", nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := Unmarshal(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
			}
		})
	}
}

func TestToBencodeDictionary(t *testing.T) {
	type TestStruct struct {
		Foo  string `bencode:"foo"`
		Bar  int    `bencode:"bar"`
		Baz  string
		Quux float64 `bencode:"quux"`
	}

	testStruct := TestStruct{
		Foo:  "hello",
		Bar:  42,
		Baz:  "world",
		Quux: 3.14,
	}

	result, err := ToBencodeDictionary(testStruct)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"foo":  "hello",
		"bar":  42,
		"Baz":  "world",
		"quux": 3.14,
	}

	assert.Equal(t, expected, result)
}
