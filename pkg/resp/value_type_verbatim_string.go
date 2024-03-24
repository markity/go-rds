package resp

import (
	"bytes"
	"errors"
	"fmt"
)

type VerbatimString struct {
	Encoding string
	Data     string
}

func (vs *VerbatimString) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if len(vs.Encoding) != 3 {
		return nil, errors.New("encoding must be exactly 3 bytes")
	}
	buf.WriteByte(':')
	buf.WriteString(vs.Data)
	return buf.Bytes(), nil
}

func (vs *VerbatimString) String() string {
	return fmt.Sprintf("verbatim string encoding: %v, data: %v", vs.Encoding, vs.Data)
}

func ToVerbatimString(encoding string, data string) (*VerbatimString, error) {
	if len(encoding) != 3 {
		return nil, errors.New("encoding must be exactly 3 bytes")
	}
	return &VerbatimString{
		Encoding: encoding,
		Data:     data,
	}, nil
}
