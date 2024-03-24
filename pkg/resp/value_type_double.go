package resp

import (
	"bytes"
	"errors"
	"fmt"
)

type Double struct {
	IsPositiveInf bool
	IsNegativeInf bool
	IsNan         bool
	Data          float64
}

func (d *Double) EncodeToBytes() ([]byte, error) {
	if d.IsNegativeInf {
		if d.IsPositiveInf || d.IsNan {
			return nil, errors.New("invalid Double")
		}
		return []byte("-inf\r\n"), nil
	}
	if d.IsPositiveInf {
		if d.IsNegativeInf || d.IsNan {
			return nil, errors.New("invalid Double")
		}
	}
	if d.IsNan {
		if d.IsPositiveInf || d.IsNegativeInf {
			return nil, errors.New("invalid Double")
		}
	}

	buf := bytes.Buffer{}
	buf.WriteByte(',')
	buf.WriteString(fmt.Sprint(d.Data))
	buf.WriteString("\r\n")
	return buf.Bytes(), nil
}

func (d *Double) String() string {
	return fmt.Sprintf("double data: %v, +inf: %v, -inf: %v, nan: %v",
		d.Data, d.IsPositiveInf, d.IsNegativeInf, d.IsNan)
}

func ToDouble(typ string, data float64) (*Double, error) {
	switch typ {
	case "nan":
		return &Double{
			IsNan: true,
		}, nil
	case "+inf":
		return &Double{
			IsPositiveInf: true,
		}, nil
	case "-inf":
		return &Double{
			IsNegativeInf: true,
		}, nil
	case "value":
		return &Double{
			Data: data,
		}, nil
	default:
		return nil, errors.New("invalid type")
	}
}
