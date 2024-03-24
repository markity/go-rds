package resp

import (
	"bytes"
	"fmt"
	"math/big"
)

type BigNum struct {
	Data big.Int
}

func (b *BigNum) EncodeToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('(')
	buf.WriteString(b.Data.Text(10))
	buf.WriteString("\r\n")
	return buf.Bytes(), nil
}

func (b *BigNum) String() string {
	return fmt.Sprintf("big num data: %v", b.Data.String())
}

func ToBigNum(i big.Int) *BigNum {
	return &BigNum{
		Data: i,
	}
}
