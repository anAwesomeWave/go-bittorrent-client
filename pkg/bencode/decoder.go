package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Decoder struct {
}

// res can be one of: int, string, []any, map[string]any, nil in case of an error
func (d Decoder) Decode(data io.Reader) (res any, err error) {
	br, ok := data.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(data)
	}
	return d.parseData(br)
}

func (d Decoder) parseData(data *bufio.Reader) (any, error) {
	// recursive. get one byte per time and process
	newByte, err := data.ReadByte()
	if err != nil {
		return nil, err
	}
	switch newByte {
	case 'i':
		// int
		return d.parseInt(data)
	case 'l':
		// list
		l := make([]any, 0, 10)
		for {
			newByte, err := data.ReadByte()
			if err != nil {
				return nil, err
			}
			if newByte == 'e' {
				break
			}
			if err := data.UnreadByte(); err != nil {
				return nil, err
			}
			elem, err := d.parseData(data)
			if err != nil {
				return nil, err
			}
			l = append(l, elem)
		}
		return l, nil
	case 'd':
		dict := make(map[string]any, 10)
		for {
			newByte, err := data.ReadByte()
			if err != nil {
				return nil, err
			}
			if newByte == 'e' {
				break
			}
			if err := data.UnreadByte(); err != nil {
				return nil, err
			}
			k, err := d.parseData(data)
			if err != nil {
				return nil, err
			}
			key, ok := k.(string)
			if !ok {
				return nil, errors.New("bencode: non-string dictionary key")
			}
			val, err := d.parseData(data)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}
		return dict, nil
	default:
		data.UnreadByte()
		stringLenBuf, err := d.readUntil(data, ':')
		if err != nil {
			return nil, err
		}
		stringLen, err := strconv.Atoi(stringLenBuf)
		s := make([]byte, stringLen)

		if _, err := io.ReadAtLeast(data, s, stringLen); err != nil {
			return nil, err
		}
		return string(s), nil
	}
}

func (d Decoder) parseInt(data *bufio.Reader) (int, error) {
	res, err := d.readUntil(data, 'e')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (d Decoder) readUntil(data *bufio.Reader, delim byte) (string, error) {
	intBuf := strings.Builder{}
	intBuf.Grow(10)
	for {
		newByte, err := data.ReadByte()
		if err != nil {
			return "", err
		}
		if newByte == delim {
			break
		}
		intBuf.WriteByte(newByte)
	}
	return intBuf.String(), nil
}
