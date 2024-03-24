package resp

import (
	"bytes"
	"fmt"
)

// Both map keys and values can be any of RESP's types.
type Map struct {
	Data map[Value]Value
}

func (m *Map) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('%')
	l := len(m.Data)
	buf.WriteString(fmt.Sprint(l))
	buf.WriteString("\r\n")
	for k, v := range m.Data {
		kE, err := k.EncodeToBytes()
		if err != nil {
			return nil, err
		}
		vE, err := v.EncodeToBytes()
		if err != nil {
			return nil, err
		}
		buf.Write(kE)
		buf.Write(vE)
	}
	return buf.Bytes(), nil
}

func (m *Map) String() string {
	return fmt.Sprintf("map data: %v", m.Data)
}

func ToMap(data map[Value]Value) *Map {
	return &Map{
		Data: data,
	}
}
