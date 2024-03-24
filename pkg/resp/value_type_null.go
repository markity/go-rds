package resp

type Null struct{}

func (*Null) EncodeToBytes() ([]byte, error) {
	return []byte("_\r\n"), nil
}

func (*Null) String() string {
	return "null"
}

// FIXME: it is not necessary
func ToNull() *Null {
	return &Null{}
}
