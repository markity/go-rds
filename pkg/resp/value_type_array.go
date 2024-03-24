package resp

import (
	"bytes"
	"fmt"
)

type Array struct {
	Data []Value
}

func (arr *Array) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('*')
	l := len(arr.Data)
	buf.WriteString(fmt.Sprint(l))
	buf.WriteString("\r\n")
	if l == 0 {
		return buf.Bytes(), nil
	}
	for _, v := range arr.Data {
		bs, err := v.EncodeToBytes()
		if err != nil {
			return nil, err
		}
		buf.Write(bs)
	}
	return buf.Bytes(), nil
}

func (arr *Array) String() string {
	return fmt.Sprintf("array data: %v", arr.Data)
}

func ToArray(vals ...Value) *Array {
	return &Array{
		Data: vals,
	}
}
