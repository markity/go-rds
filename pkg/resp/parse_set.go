package resp

type parseSetStateEnum int

const (
	parseSetWantingLength parseSetStateEnum = iota
	parseSetWantingData
)

type parseSetState struct {
	state  parseSetStateEnum
	length int
	data   []Value
}

func (st *parseSetState) SetEntry(v Value) {
	st.data = append(st.data, v)
}
