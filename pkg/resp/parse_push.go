package resp

type parsePushStateEnum int

const (
	parsePushWantingLength parsePushStateEnum = iota
	parsePushWantingData
)

type parsePushState struct {
	state  parsePushStateEnum
	length int
	data   []Value
}

func (st *parsePushState) PushEntry(v Value) {
	st.data = append(st.data, v)
}
