package resp

import (
	"bytes"
	"fmt"
)

// This type is a CRLF-terminated string that represents a signed, base-10, 64-bit integer.

type Integer struct {
	Data int64
}

func (i *Integer) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(':')
	buf.WriteString(fmt.Sprint(i.Data))
	buf.WriteString("\r\n")
	fmt.Println(buf.Bytes())
	return buf.Bytes(), nil
}

func (i *Integer) String() string {
	return fmt.Sprintf("integer data: %v", i.Data)
}

func ToInteger(data int64) *Integer {
	return &Integer{Data: data}
}
