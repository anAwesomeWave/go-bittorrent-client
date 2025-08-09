package bencode

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

type testData struct {
	testName    string
	encoded     string
	expected    any
	expectedErr error
}

func TestDecodeSimpleUseCases(t *testing.T) {
	data := []testData{
		{
			"Check empty decoding",
			"",
			nil,
			io.EOF,
		},
		{
			"Check integer decoding",
			"i42e",
			42,
			nil,
		},
		{
			"Check string decoding",
			"5:abcde",
			"abcde",
			nil,
		},
	}

	for _, testCase := range data {
		encodedData := testCase.encoded
		expectedDecoded := testCase.expected

		actualDecoded, err := DecodeString(encodedData)
		t.Logf("Test: %s\n", testCase.testName)
		if !errors.Is(err, testCase.expectedErr) {
			t.Errorf("Expected err %v.\tGot: %v", testCase.expectedErr, err)
		}
		if actualDecoded != expectedDecoded {
			t.Errorf("Expected %v type %T, got %v, type %T", expectedDecoded, expectedDecoded, actualDecoded, actualDecoded)
		}
	}
}

func TestDecodeDeepStructures(t *testing.T) {
	data := []testData{
		{
			"Check list of any decoding",
			"li42e5:abcdee",
			[]any{42, "abcde"},
			nil,
		},
		{
			"Check nested lists of any decoding",
			"lli1el1:1eei42e5:abcdee",
			[]any{[]any{1, []any{"1"}}, 42, "abcde"},
			nil,
		},
		{
			"Check maps of any decoding",
			"d7:meaningi42e4:wiki7:bencodee",
			map[string]any{"meaning": 42, "wiki": "bencode"},
			nil,
		},
		{
			"Check nested maps of any decoding",
			"d7:meaningd1:ai1ee4:wiki7:bencodee",
			map[string]any{"meaning": map[string]any{"a": 1}, "wiki": "bencode"},
			nil,
		},
		{
			"Check nested lists/maps of any decoding",
			"ld7:meaningd1:ai1ee4:wiki7:bencodeed1:ali1eeee",
			[]any{map[string]any{"meaning": map[string]any{"a": 1}, "wiki": "bencode"}, map[string]any{"a": []any{1}}},
			nil,
		},
	}

	for _, testCase := range data {
		encodedData := testCase.encoded
		expectedDecoded := testCase.expected

		actualDecoded, err := DecodeString(encodedData)
		t.Logf("Test: %s\n", testCase.testName)
		if !errors.Is(err, testCase.expectedErr) {
			t.Errorf("Expected err %v.\tGot: %v", testCase.expectedErr, err)
		}
		if !reflect.DeepEqual(expectedDecoded, actualDecoded) {
			t.Errorf("Expected %v type %T, got %v, type %T", expectedDecoded, expectedDecoded, actualDecoded, actualDecoded)
		}
	}
}
