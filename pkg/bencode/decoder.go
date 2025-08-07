package bencode

import (
	"bufio"
	"io"
)

type Decoder struct {
}

// res can be one of: int, string, []any, map[string]any
func (d Decoder) Decode(data io.Reader) (res any, err error) {
	br, ok := data.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(data)
	}
	return d.parseData(br)
}

func (d Decoder) parseData(data *bufio.Reader) (any, error) {
	// recursive. get one byte per time and process
	return nil, nil
}
