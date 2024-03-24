package resp

import (
	"bytes"
	"fmt"
)

type Bool struct {
	Data bool
}

func (b *Bool) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('#')
	if b.Data {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteString("\r\n")
	return buf.Bytes(), nil
}

func (b *Bool) String() string {
	return fmt.Sprintf("bool data: %v", b.Data)
}

func ToBool(data bool) *Bool {
	return &Bool{
		Data: data,
	}
}
