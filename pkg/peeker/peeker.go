package peeker

type Peeker struct {
	data []byte
}

func (pk *Peeker) Peek() []byte {
	return pk.data
}

func (pk *Peeker) LenRemain() int {
	return int(len(pk.data))
}

func (pk *Peeker) Write(bs []byte) {
	pk.data = append(pk.data, bs...)
}

func (pk *Peeker) Discard(n int) {
	pk.data = pk.data[n:]
}
