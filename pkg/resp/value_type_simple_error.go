package resp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type SimpleError struct {
	ErrType string
	ErrMsg  string
}

func (s *SimpleError) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('-')
	buf.WriteString(s.ErrType)
	buf.WriteByte(' ')
	buf.WriteString(s.ErrMsg)
	buf.WriteString("\r\n")
	return buf.Bytes(), nil
}

func (s *SimpleError) String() string {
	return fmt.Sprintf("simple error err_type: %v, err_msg: %v", s.ErrType, s.ErrMsg)
}

func ToSimpleError(errtype string, errmsg string) (*SimpleError, error) {
	for _, v := range errtype {
		if !unicode.IsUpper(v) {
			return nil, errors.New("err type must be upper case letters")
		}
	}
	if strings.Contains(errtype, "\r") {
		return nil, errors.New("cannot contain \\r")
	}
	if strings.Contains(errtype, "\n") {
		return nil, errors.New("cannot contain \\n")
	}
	if strings.Contains(errmsg, "\r") {
		return nil, errors.New("cannot contain \\r")
	}
	if strings.Contains(errmsg, "\n") {
		return nil, errors.New("cannot contain \\n")
	}

	return &SimpleError{ErrType: errtype, ErrMsg: errmsg}, nil
}
