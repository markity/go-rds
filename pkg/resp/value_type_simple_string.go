package resp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type SimpleString struct {
	Data string
}

func (s *SimpleString) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('+')
	buf.WriteString(s.Data)
	buf.WriteString("\r\n")
	return buf.Bytes(), nil
}

func (s *SimpleString) String() string {
	return fmt.Sprintf("simple string data: %v", s.Data)
}

func ToSimpleString(data string) (*SimpleString, error) {
	/*
		From: https://redis.io/docs/reference/protocol-spec/#simple-strings
		The string mustn't contain a CR (\r) or LF (\n) character
	*/
	if strings.Contains(data, "\r") {
		return nil, errors.New("cannot contain \\r")
	}
	if strings.Contains(data, "\n") {
		return nil, errors.New("cannot contain \\n")
	}

	return &SimpleString{Data: data}, nil
}
