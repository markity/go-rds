package resp

type parseArrayStateEnum int

const (
	parseArrayWantingLength parseArrayStateEnum = iota
	parseArrayWantingData
)

type parseArrayState struct {
	state  parseArrayStateEnum
	length int
	data   []Value
}

func (st *parseArrayState) PushEntry(v Value) {
	st.data = append(st.data, v)
}
