package resp

type Value interface {
	EncodeToBytes() ([]byte, error)
	String() string
}
