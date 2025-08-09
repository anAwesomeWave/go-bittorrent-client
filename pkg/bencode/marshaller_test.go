package bencode

import (
	"errors"
	"reflect"
	"testing"
)

type Decoded struct {
	ExampleString string `mapstructure:"example_string"`
	Str           string
	ExampleInt    int `mapstructure:"example_int"`
	I             int
	Data          []string
	Extra         map[string]int
}

type marshallerTestData struct {
	testName       string
	encoded        string
	expectedValues Decoded
	expectedErr    error
}

func TestUnmarshal(t *testing.T) {
	data := []marshallerTestData{
		{
			"Check struct filling",
			"d14:example_string5:test13:str5:test211:example_inti1e1:ii2e4:datal1:a1:b1:ce5:extrad1:ai1eee",
			Decoded{
				ExampleString: "test1",
				Str:           "test2",
				ExampleInt:    1,
				I:             2,
				Data:          []string{"a", "b", "c"},
				Extra:         map[string]int{"a": 1},
			},
			nil,
		},
	}

	for _, testCase := range data {
		encodedData := testCase.encoded
		expectedDecoded := testCase.expectedValues
		actualDecoded := Decoded{}

		err := UnmarshalString(encodedData, &actualDecoded)
		t.Logf("Test: %s\n", testCase.testName)
		if !errors.Is(err, testCase.expectedErr) {
			t.Errorf("Expected err %v.\tGot: %v", testCase.expectedErr, err)
		}
		if !reflect.DeepEqual(expectedDecoded, actualDecoded) {
			t.Errorf("Expected %v type %T, got %v, type %T", expectedDecoded, expectedDecoded, actualDecoded, actualDecoded)
		}
	}
}
