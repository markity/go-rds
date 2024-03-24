package resp

import (
	"bytes"
	"fmt"
)

// The empty string's encoding is: $0\r\n\r\n
type BulkString struct {
	Data string
}

func (bs *BulkString) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('$')
	buf.WriteString(fmt.Sprint(len(bs.Data)))
	buf.WriteString("\r\n")
	buf.WriteString(bs.Data)
	buf.WriteString("\r\n")
	return buf.Bytes(), nil
}

func (bs *BulkString) String() string {
	return fmt.Sprintf("bulk string data: %v", bs.Data)
}

func ToBulkString(data string) *BulkString {
	return &BulkString{
		Data: data,
	}
}
