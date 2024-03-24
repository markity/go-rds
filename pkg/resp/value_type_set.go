package resp

import (
	"bytes"
	"fmt"
)

type Set struct {
	Data []Value
}

func (arr *Set) EncodeToBytes() ([]byte, error) {
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

func (arr *Set) String() string {
	return fmt.Sprintf("set data: %v", arr.Data)
}

func ToSet(vals ...Value) *Set {
	return &Set{
		Data: vals,
	}
}
