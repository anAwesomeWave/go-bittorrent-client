package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func DecodeString(s string) (any, error) {
	return Decode(strings.NewReader(s))
}

// Decode res can be one of: int, string, []any, map[string]any, nil may be returned in case of an error
func Decode(data io.Reader) (res any, err error) {
	br, ok := data.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(data)
	}
	return parseData(br)
}

func parseData(data *bufio.Reader) (any, error) {
	// recursive. get one byte per time and process
	newByte, err := data.ReadByte()
	if err != nil {
		return nil, err
	}
	switch newByte {
	case 'i':
		// int
		return parseInt(data)
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
			elem, err := parseData(data)
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
			k, err := parseData(data)
			if err != nil {
				return nil, err
			}
			key, ok := k.(string)
			if !ok {
				return nil, errors.New("bencode: non-string dictionary key")
			}
			val, err := parseData(data)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}
		return dict, nil
	default:
		if err := data.UnreadByte(); err != nil {
			return nil, err
		}
		return parseString(data)
	}
}

func parseString(data *bufio.Reader) (string, error) {
	stringLenBuf, err := readUntil(data, ':')
	if err != nil {
		return "", err
	}
	stringLen, err := strconv.Atoi(stringLenBuf)
	s := make([]byte, stringLen)

	if _, err := io.ReadAtLeast(data, s, stringLen); err != nil {
		return "", err
	}
	return string(s), nil
}

func parseInt(data *bufio.Reader) (int, error) {
	res, err := readUntil(data, 'e')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func readUntil(data *bufio.Reader, delim byte) (string, error) {
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
