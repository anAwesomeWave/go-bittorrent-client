package bencode

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"reflect"
	"strings"
)

// marshaller takes bencoded data as input and stores as struct. it is slow because uses decoder's map[string]any as
// intermediate step

func UnmarshalString(s string, v interface{}) error {
	return Unmarshal(strings.NewReader(s), v)
}

func Unmarshal(data io.Reader, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr || v == nil {
		return fmt.Errorf("v must be a non-nil pointer to a struct")
	}
	prepared, err := Decode(data)
	if err != nil {
		return err
	}

	return mapstructure.Decode(prepared, v)
}
